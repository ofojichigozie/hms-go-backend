# Hospital Management System

A comprehensive hospital management system built with Go, Gin, GORM, and PostgreSQL. This system provides APIs for managing staff, patients, appointments, and clinical notes with role-based access control.

## Features

- **Authentication & Authorization**: Secure login and role-based access control
- **Staff Management**: Add, view, update, and delete staff members
- **Patient Management**: Register, update, and manage patient information
- **Appointment Scheduling**: Create and manage patient appointments
- **Clinical Notes**: Create and manage clinical notes for patient visits
- **Role-Based Access**: Different permissions for Admins, Doctors, and Receptionists

## Live Demo

Access the deployed API at: [https://hms-go-api.onrender.com](https://hms-go-api.onrender.com)

## Getting Started

### Prerequisites

- Go 1.16+
- PostgreSQL
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/ofojichigozie/hms-go-backend.git
   cd hms-go-backend
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up environment variables:
   Create a `.env` file in the root directory with the following variables:
   ```
   APP_ENV=development
   PORT=5000
   DB_URL=postgres://username:password@localhost:5432/hms_db
   JWT_SECRET=your_jwt_secret_key
   ```

4. Set up the database:
   ```bash
   createdb hms_db   # If using PostgreSQL CLI
   ```

5. Run database migrations:
   ```bash
   go run migrate/migrate.go

### Running Locally

#### With Air (Hot Reload)

The project is configured with Air for hot reloading during development.

```bash
# Install Air if you haven't already
go install github.com/cosmtrek/air@latest

# Run the application with Air
air
```

#### Without Air

```bash
go run main.go
```

The server will start on `http://localhost:5000`.

## API Endpoints

### Authentication
- `POST /auth/login` - User login
- `POST /auth/refresh` - Refresh authentication token (requires authentication)

### Staff Management
- `POST /staff` - Create new staff (Admin only)
- `GET /staff` - Get all staff (Admin only)
- `GET /staff/:id` - Get staff by ID (Authenticated users)
- `PATCH /staff/:id` - Update staff (Admin only)
- `DELETE /staff/:id` - Delete staff (Admin only)

### Patient Management
- `POST /patients` - Register new patient (Receptionist only)
- `GET /patients` - Get all patients (Receptionist and Doctor)
- `GET /patients/:id` - Get patient by ID (Receptionist and Doctor)
- `PATCH /patients/:id` - Update patient (Receptionist only)
- `DELETE /patients/:id` - Delete patient (Receptionist only)

### Appointment Management
- `POST /appointments` - Create appointment (Receptionist only)
- `GET /appointments` - Get all appointments (Receptionist and Doctor)
- `GET /appointments/:id` - Get appointment by ID (Receptionist and Doctor)
- `PATCH /appointments/:id` - Update appointment (Receptionist only)
- `DELETE /appointments/:id` - Delete appointment (Receptionist only)

### Clinical Notes
- `POST /clinical-notes` - Create clinical note (Doctor only)
- `GET /clinical-notes/:id` - Get note by ID (Doctor and Receptionist)
- `GET /clinical-notes/patient/:patientId` - Get notes by patient ID (Doctor and Receptionist)
- `PATCH /clinical-notes/:id` - Update clinical note (Doctor only)
- `DELETE /clinical-notes/:id` - Delete clinical note (Doctor only)

## Project Structure

```
hms-go-backend/
├── constants/         # Application constants
├── controllers/       # Request handlers
├── initializers/      # Database setup and configuration
├── middleware/        # Authentication and authorization middleware
├── models/            # Database models
├── repositories/      # Database interaction logic
├── routes/            # API route definitions
├── services/          # Business logic
├── utils              # Utility functions
├── .air.toml          # Air configuration for hot reloading
├── .env               # Environment variables
├── go.mod             # Go module definition
├── go.sum             # Go module checksums
└── main.go            # Application entry point
```

## Testing

Run tests with:

```bash
go test ./services
```

For test coverage:

```bash
go test ./services -cover
```

## Initial Setup

On first run, the system automatically creates an initial admin user. You can use this account to create other staff members.

## License

[MIT License](LICENSE)

## Author

Ofojichigozie