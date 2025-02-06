import httpx
import asyncio
import logging
import json
import gdown
import pandas as pd
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
    'gemini-2.0-flash',
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

# Create a matching agent to map website team names to correct names.
match_agent = Agent(
    model,
    system_prompt=(
        "You are a name matcher. Given two lists of team names, one from a website and one from an external source, "
        "map each website team name to its best match from the external list. "
        "Return only plain JSON (no markdown) in the format "
        "{\"WebsiteTeamName\": \"CorrectTeamName\", ...}."
    ),
    deps_type=None,
    result_type=str
)


async def match_team_names(website_names: List[str], correct_names: List[str]) -> Dict[str, str]:
    prompt = (
            "Website team names: " + json.dumps(website_names) + "\n"
                                                                 "Correct team names: " + json.dumps(
        correct_names) + "\n"
                         "Map each website team name to its best match from the correct names list. "
                         "Return only plain JSON (no markdown) in the format {\"WebsiteTeamName\": \"CorrectTeamName\", ...}."
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
    # Read the entire sheet without assuming any header.
    df = pd.read_excel(xlsx_path, sheet_name="ShowProgram", header=None)
    # Flatten the DataFrame to extract all text cells.
    names = df.astype(str).values.flatten().tolist()
    # Remove empty strings and duplicates.
    names = [n.strip() for n in names if n.strip()]
    return list(dict.fromkeys(names))  # preserve order and remove duplicates


# -------------------------------
# Main Crawling Logic
# -------------------------------

async def main():
    base_url = "https://improneuf.com"
    overview_url = base_url + "/dt/web/teams-overview"
    try:
        # Dummy RunContext (no dependencies)
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
        # Fetch teams overview HTML.
        overview_html = await fetch_html(dummy_ctx, url=overview_url)
        teams = await extract_teams(overview_html)

        # Download Excel file from Google Drive and extract all text.
        xlsx_file = download_excel_file()
        correct_names = extract_correct_names(xlsx_file)

        website_names = [team.name for team in teams]
        # Use the matching agent to create a mapping.
        name_mapping = await match_team_names(website_names, correct_names)

        # Replace website team names with the correct names if available.
        for team in teams:
            if team.name in name_mapping:
                team.name = name_mapping[team.name]

        # Crawl each team's page and extract the team image.
        team_images: Dict[str, str] = {}
        for team in teams:
            team_page_url = base_url + team.link
            team_html = await fetch_html(dummy_ctx, url=team_page_url)
            image_url = extract_team_image(team_html)
            team_images[team.name] = base_url + image_url

        # Print the final mapping as JSON.
        with open("team-images.json", "w", encoding="utf-8") as f:
            json.dump(team_images, f, indent=2, ensure_ascii=False)
    except Exception as e:
        logging.error(f"An error occurred: {e}")


if __name__ == "__main__":
    asyncio.run(main())
