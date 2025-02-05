import argparse
import os
import json
import requests
from io import BytesIO
from PIL import Image

def main():
    parser = argparse.ArgumentParser(
        description="Download images from URLs specified in a JSON file, convert them to PNG, and save them."
    )
    parser.add_argument("input_file", help="Path to the JSON file mapping team names to image URLs")
    parser.add_argument("output_dir", help="Directory where the PNG images will be saved")
    args = parser.parse_args()

    # Load the JSON file that maps team names to image URLs.
    with open(args.input_file, "r", encoding="utf-8") as f:
        team_images = json.load(f)

    # Create the output directory (overwrite files if they exist).
    os.makedirs(args.output_dir, exist_ok=True)

    # Define characters that are invalid in Linux filenames.
    invalid_chars = {"/", "\0"}

    for team, url in team_images.items():
        # Skip team names that contain any invalid characters.
        if any(ch in team for ch in invalid_chars):
            print(f"Skipping team '{team}' due to invalid characters in its name.")
            continue

        try:
            response = requests.get(url)
            response.raise_for_status()
            # Load the image from the response.
            img = Image.open(BytesIO(response.content))
            # Convert the image to PNG.
            img_converted = img.convert("RGB")
            filename = team + ".png"
            path = os.path.join(args.output_dir, filename)
            img_converted.save(path, "PNG")
            print(f"Saved {path}")
        except Exception as e:
            print(f"Error processing {team}: {e}")

if __name__ == "__main__":
    main()
