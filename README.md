# management

https://improneuf.github.io/management/


# Generate workshop banners

# Step 1: Generate HTML files in open-workshops root
cd src/open-workshops
go run main.go
# Creates: workshop-1-Emotional-Whirlwinds.html, workshop-2-Building-a-Universe.html, etc.

# Step 2: Convert to PNG and save in output directory
cd ../html-to-png
go run convert-workshops.go
# Creates: open-workshops/output/workshop-1-Emotional-Whirlwinds-fb.jpg, etc.