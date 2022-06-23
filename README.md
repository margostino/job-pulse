# JOB Pulse

JobPulse is a hiring market analyzer application to answer the following questions: 
  
- How's hiring of company X going?
- Is there any trending role?
- Is there any visible hiring freeze? First clue of Layoff?

JobPulse searches and collects Job Posts from any source (e.g. Linkedin) for a given Position (e.g. Software Engineer) and Location (e.g. Stockholm) and ingests into a [MongoDB Atlas Database on Cloud](https://www.mongodb.com/cloud/atlas). 

## Architecture

This project uses Golang and MongoDB Atlas on Cloud. 
So far only a data collector application is implemented.
The Data collector runs automatically **once a day**. 

![Architecture](./assets/images/architecture.png#center)

## Usage

### Dashboard:

If you want to explore the public charts, you can access [here](https://charts.mongodb.com/charts-project-0-mcjod/public/dashboards/62ab5b86-5868-44fe-885f-14caf30ccad1).

![MongoDB Charts](./assets/images/mongodb-charts.png#center)

### Collector:

If you want to run the data collector by your own:
```go
go run ./runner "software engineer" "stockholm"
```

## Features
TBD

## Contribution
TBD

## Brainstorming
1. Automate collection: Github action ⏱ -> Vercel function
2. Make charts/dashboards public
3. CLI to query data
4. Multi source integration
5. Event Correlation (?)
6. Reporting automation
7. Improve logging
8. Alarms for a given rule (hiring freeze?)
9. Testing, testing, testing
10. Normalize and improve text sanitization
11. Geo Chart
12. IO async
13. Support batches or streams

