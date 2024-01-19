# scraping-keyword-web
Keyword Scraping Web Apps
Welcome to our keyword scraping web app repository! This repository contains both the backend (Golang) and frontend (ReactJS) applications. Our app is designed to collect information about search keywords, providing a valuable tool for gathering and analyzing keyword data.

Starting the App
To start the app, follow the steps below:

1. Backend (Golang)
Start the Server: Use the following command to start the PosgreSQL DB and Golang server:
yarn start-api

2. Frontend (ReactJS)
Start the Development Server: Use the following command to start the development server for the frontend application:
npm start
Running Both Backend and Frontend Together

To run both the backend and frontend together, you can use the following command:
docker-compose up -d && npm install && go run ./backend/main.go
Feel free to reach out if you have any further questions or need additional assistance!