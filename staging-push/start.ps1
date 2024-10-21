# Start Air live reload for Go server
Start-Process -FilePath "air" -WorkingDirectory "C:\Users\menzi\go\code\mango-monopoly\v1"

# Start Tailwind CSS watch process
Start-Process -FilePath "tailwindcss" -ArgumentList "-i .\ui\static\css\input.css -o .\ui\static\css\output.css --watch" -WorkingDirectory "C:\Users\menzi\go\code\mango-monopoly\v1"
