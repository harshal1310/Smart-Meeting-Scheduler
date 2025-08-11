# Smart Scheduler

A sophisticated meeting scheduling API built with Go that intelligently finds optimal meeting times for multiple participants while considering their existing calendar events and applying smart heuristics for the best user experience.

### Core Algorithm: Optimal Meeting Time Detection

The Smart Scheduler uses a sophisticated multi-step algorithm to find the best meeting times:


### Heuristics Implemented

1. ** Time Preference**: Morning slots (9-12 PM) scored higher than afternoon
2. ** Buffer Management**: 15-minute buffers preferred around meetings
3. ** Gap Optimization**: Avoids creating unusable 30-minute gaps
4. ** Business Hours**: Penalties for scheduling outside 9 AM - 4 PM
5. ** Conflict Avoidance**: Absolute prevention of overlapping meetings
6. ** Multiple Participants**: Considers all attendees' calendars simultaneously


### Local Development Setup

#### 1. Clone the Repository
```bash
git clone https://github.com/harshal1310/Smart-Meeting-Scheduler
cd Smart-Scheduler
```

#### 2. Install Dependencies
```bash
go mod download
```

#### 3. Database Setup

**Option A: Local PostgreSQL**
```bash
# Create database
createdb smartscheduler

# Set environment variables
export DATABASE_URL="postgres://username:password@localhost/smart_scheduler?sslmode=disable"
export DB_NAME="smartscheduler"
export PORT="8080"
```

**Option B: Docker PostgreSQL**
```bash
# Run PostgreSQL in Docker
docker run --name postgres-scheduler \
  -e POSTGRES_DB=smart_scheduler \
  -e POSTGRES_USER=scheduler \
  -e POSTGRES_PASSWORD=password \
  -p 5432:5432 \
  -d postgres:15

# Set environment variables
export DATABASE_URL="postgres://scheduler:password@localhost/smart_scheduler?sslmode=disable"
export DB_NAME="smart_scheduler"
export PORT="8080"
```

#### 4. Build & Run

**Development Mode:**
```bash
# Run with automatic recompilation
go run cmd/server/main.go
```

**Production Build:**
```bash
# Build binary
go build -o smart-scheduler cmd/server/main.go

# Run binary
./smart-scheduler
```

#### 5. Verify Installation
```bash
# Check server health
curl http://localhost:8080/api/v1/calendar/user1

# Expected: JSON response with user's calendar events
```

##  API Endpoints

### Base URL
```
http://localhost:8080/api/v1
```

### Complete Endpoint URLs
- **POST** `http://localhost:8080/api/v1/schedule` - Schedule a new meeting
- **GET** `http://localhost:8080/api/v1/calendar/{userID}` - Get user's calendar events

### Endpoints

#### 1. **Schedule Meeting**
```http
POST /api/v1/schedule
Content-Type: application/json

{
  "title": "Team Standup",
  "userIDs": ["user1", "user2", "user3"],
  "durationMinutes": 60,
  "timeRange": {
    "start": "2025-08-09T09:00:00+05:30",
    "end": "2025-08-09T17:00:00+05:30"
  }
}
```

**Response:**
```json
{
  "meetingId": "meeting-12345",
  "title": "Team Standup",
  "participantIds": ["user1", "user2", "user3"],
  "startTime": "2025-08-09T10:00:00+05:30",
  "endTime": "2025-08-09T11:00:00+05:30"
}
```

#### 2. **Get User Calendar**
```http
GET /api/v1/calendar/{userID}?start=2025-08-09T08:00:00+05:30&end=2025-08-09T18:00:00+05:30
```

**Response:**
```json
{
  "events": [
    {
      "id": 1,
      "eventCode": "event1",
      "userId": "user1",
      "title": "Daily Standup",
      "startTime": "2025-08-09T09:00:00+05:30",
      "endTime": "2025-08-09T09:30:00+05:30"
    }
  ]
}
```

### Error Responses

```json
{
  "error": "Error message description",
  "code": 400
}
```

**Common Status Codes:**
- `200`: Success
- `400`: Bad Request (invalid JSON, missing fields)
- `404`: User not found
- `409`: Conflict (no available time slots found)
- `500`: Internal server error

##  Testing

### Running Tests

#### **All Tests**
```bash
go test ./... -v
```

#### **Package-Specific Tests**
```bash
# API utilities
go test ./api -v

# Data models
go test ./model -v

# Repository layer
go test ./repository -v

# Business logic
go test ./service -v

# HTTP handlers
go test ./handlers -v
```

#### **With Coverage**
```bash
# Generate coverage report
go test ./... -coverprofile=coverage.out

# View coverage in browser
go tool cover -html=coverage.out
```

#### **Integration Tests**
```bash
# Run main integration tests
go test ./main_test.go -v
```
