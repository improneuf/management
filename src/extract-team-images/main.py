import httpx
import asyncio
import logging
import json
import gdown
import difflib
from bs4 import BeautifulSoup
from pydantic import BaseModel
from pydantic_ai import Agent, RunContext
from typing import List, Dict
from pydantic_ai.models.vertexai import VertexAIModel

# Define the Pydantic model for teams (from the website).
class Team(BaseModel):
    name: str
    link: str

# Helper function to remove markdown code block delimiters.
def clean_markdown_json(s: str) -> str:
    lines = s.strip().splitlines()
    if lines and lines[0].startswith("```"):
        lines = lines[1:]
    if lines and lines[-1].startswith("```"):
        lines = lines[:-1]
    return "\n".join(lines)

# Helper function to extract the team image URL from a team page.
def extract_team_image(html: str) -> str:
    soup = BeautifulSoup(html, "html.parser")
    div = soup.find("div", class_=lambda c: c and "field--name-field-team-photo-for-promo-imgno" in c)
    if div:
        img = div.find("img")
        if img and img.get("src"):
            return img["src"]
    return ""

# Initialize the VertexAI model and extraction agent.
model = VertexAIModel(
    'gemini-2.0-flash-lite-preview-02-05',
    project_id="geminitest-407623",
    service_account_file="geminitest-407623-91eaab2a74e8.json"
)
extraction_agent = Agent(
    model,
    system_prompt=(
        "Extract all team names and their corresponding links from the provided HTML. "
        "Return the result as plain JSON (no markdown formatting) in the form of a JSON array "
        "of objects with keys 'name' and 'link'."
    ),
    deps_type=None,
    result_type=str  # We'll clean and parse the string manually.
)

# Define a tool to fetch HTML content.
@extraction_agent.tool
async def fetch_html(ctx: RunContext[None], url: str) -> str:
    async with httpx.AsyncClient() as client:
        response = await client.get(url)
        response.raise_for_status()
        return response.text

# Extraction logic: Send HTML to the agent, clean the output, and parse it into a list of Team objects.
async def extract_teams(html: str) -> List[Team]:
    prompt = (
        "Extract all team names and their corresponding links from the following HTML. "
        "Return the result as plain JSON (no markdown formatting) in the form of a JSON array "
        "of objects with keys 'name' and 'link'. HTML: " + html
    )
    result = await extraction_agent.run(prompt)
    cleaned = clean_markdown_json(result.data)
    data = json.loads(cleaned)
    teams = [Team(**item) for item in data]
    return teams

# -------------------------------
# Agent for Fuzzy Name Matching
# -------------------------------
match_agent = Agent(
    model,
    system_prompt=(
        "You are a name matcher. Given two lists of team names, one from an external source "
        "and one from a website, map each external team name to its best match from the website list. "
        "Return only plain JSON (no markdown) in the format "
        "{\"ExternalTeamName\": \"WebsiteTeamName\", ...}."
    ),
    deps_type=None,
    result_type=str
)

async def match_team_names_excel_to_website(excel_names: List[str], website_names: List[str]) -> Dict[str, str]:
    prompt = (
        "You are a name matcher. Given two lists of team names, one from an external source and one from a website, "
        "map each external team name to its best match from the website list. "
        "Return only plain JSON (no markdown) in the format "
        "{\"ExternalTeamName\": \"WebsiteTeamName\", ...}.\n"
        "External team names: " + json.dumps(excel_names) + "\n"
        "Website team names: " + json.dumps(website_names)
    )
    result = await match_agent.run(prompt)
    cleaned = clean_markdown_json(result.data)
    mapping = json.loads(cleaned)
    return mapping

# -------------------------------
# Excel Extraction from Google Drive
# -------------------------------
def download_excel_file() -> str:
    # Download the file from Google Drive using gdown.
    url = "https://drive.google.com/uc?id=1BYucz1R4IoH5whYe4goRbk_kO8LosrZ2"
    output = "showprogram.xlsx"
    gdown.download(url, output, quiet=False)
    return output


def extract_correct_names(xlsx_path: str) -> List[str]:
    import pandas as pd
    # Load the entire sheet without headers.
    df = pd.read_excel(xlsx_path, sheet_name="ShowProgram", header=None)

    # Attempt to locate the header row by checking each cell for a match.
    header_row = None
    for idx, row in df.iterrows():
        for cell in row.astype(str):
            cell_lower = cell.lower()
            # Look for a cell that seems to be the header for TimeSlot 1.
            if "timeslot 1" in cell_lower and "20" in cell_lower:
                header_row = idx
                break
        if header_row is not None:
            break

    if header_row is None:
        print("Header row not found. First few rows:")
        print(df.head(10))
        raise ValueError("Could not locate header row with the desired column names.")

    # Use the found row as the header.
    df.columns = df.iloc[header_row]
    df = df.drop(index=header_row).reset_index(drop=True)

    # Define the exact column names we expect.
    desired_cols = ["TimeSlot 1 (20 min)", "TimeSlot 2 (20 min)", "TimeSlot 3 (20 min)"]
    names = []
    for col in desired_cols:
        if col in df.columns:
            # Drop missing values, convert to string, and extend our list.
            names.extend(df[col].dropna().astype(str).tolist())

    # Clean up whitespace and remove empty strings.
    names = [n.strip() for n in names if n.strip()]
    return names


# -------------------------------
# Main Crawling Logic
# -------------------------------
async def main():
    base_url = "https://improneuf.com"
    overview_url = base_url + "/dt/web/teams-overview"
    try:
        dummy_ctx = RunContext(
            deps=None,
            model=None,
            usage=None,
            prompt="",
            messages=[],
            tool_name=None,
            retry=0,
            run_step=0
        )
        # Fetch teams overview HTML from the website.
        overview_html = await fetch_html(dummy_ctx, url=overview_url)
        teams = await extract_teams(overview_html)

        # Download Excel file and extract all team names (keeping all variations).
        xlsx_file = download_excel_file()
        excel_names = extract_correct_names(xlsx_file)

        # Build a lookup for website teams and get the list of website names.
        website_team_lookup = {team.name: team for team in teams}
        website_names = list(website_team_lookup.keys())

        # Use fuzzy matching to map each Excel team name to the best matching website team name.
        # After obtaining the fuzzy mapping:
        name_mapping = await match_team_names_excel_to_website(excel_names, website_names)

        import difflib

        # Fallback: For any Excel team with an empty mapping, try a direct case-insensitive match.
        for external in excel_names:
            if not name_mapping.get(external, "").strip():
                # Try direct match first.
                direct_match = next(
                    (wt for wt in website_names if wt.strip().lower() == external.strip().lower()),
                    None
                )
                if direct_match:
                    name_mapping[external] = direct_match
                else:
                    # If no direct match, use difflib's fuzzy matching.
                    close = difflib.get_close_matches(external, website_names, n=1, cutoff=0.8)
                    if close:
                        name_mapping[external] = close[0]

        # For each Excel team name, find the corresponding website team and fetch the image.
        team_images: Dict[str, str] = {}
        for excel_name in excel_names:
            website_match = name_mapping.get(excel_name)
            if website_match and website_match in website_team_lookup:
                team = website_team_lookup[website_match]
                team_page_url = base_url + team.link
                team_html = await fetch_html(dummy_ctx, url=team_page_url)
                image_url = extract_team_image(team_html)
                team_images[excel_name] = base_url + image_url
            else:
                logging.error(f"No website match for Excel team: {excel_name}")
                team_images[excel_name] = ""

        # Write the final mapping to JSON.
        with open("team-images.json", "w", encoding="utf-8") as f:
            json.dump(team_images, f, indent=2, ensure_ascii=False)
    except Exception as e:
        logging.error(f"An error occurred: {e}")

if __name__ == "__main__":
    asyncio.run(main())
