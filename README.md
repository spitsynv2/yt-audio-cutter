### Deploy
1. Fill `DROPBOX_TOKEN` in server.env file
2. Build image via ./build.sh
3. Run compose file

### Usage
Create job

POST http://localhost:8080/jobs
{
  "name": "Rick", 
  "youtube_url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
  "start_time": "10s",
  "end_time": "15s"
} </br>
Response: id

---

List job/bos

GET http://localhost:8080/jobs or http://localhost:8080/jobs/id </br>
Response: jobs/job

---

Process job

GET http://localhost:8080/jobs/process/id </br>
Response: dropbox download url for your cutted mp3
