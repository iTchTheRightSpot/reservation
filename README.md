# Reservation App (Angular, Go, & PostgreSQL)

A feature-rich reservation system for small businesses, like salons.

## Core Features

### 1. **Monolith Application**
   - **Frontend**: Angular served from Go binary.

### 2. **Authentication & Authorization**
   - **Login**: Email & password-based.
   - **Roles & Permissions**: Role-based access; authorized via 'GSESSION' cookie.

### 3. **Staff Management**
   - **Schedule**:
     - Only staff with `WRITE` permission can create schedules.
     - Bulk schedule creation & no conflicts.
     - Schedules can only be deleted if not linked to a reservation.
     - Weekly recurring schedules.
   - **Services**:
     - Staff with `WRITE` or `DELETE` permission can manage services.
     - Services linked to other data like reservation, or staff can’t be deleted but can be hidden.
     - Only staff with `WRITE` can assign services to others.
   - **Reservations**:
     - Allow Multiple services per reservation.
     - Quoted price required.
     - No overbooking; cannot delete, only cancel.
     - Staff notified on creation or status change.
     - Reservation status set to `COMPLETED` only with payment info.
     - 1-day reminder notifications for pending appointments.

### 4. **Client Management**
   - **Services**:
     - Clients can view all services.
     - Multiple services can be reserved.
     - Max price per service: `DECIMAL(6, 2)`.
   - **Reservations**:
     - Clients see valid available times per staff in their timezone.
     - Reservations for future dates only; no rescheduling, only cancellations.
     - Clients and staff receive notifications on status changes.

### 5. **Payment**
   - **Invoices**: Only staff can send invoices.
   - **Online Payment**: Clients can pay via an online invoice.

## Development docs
1. [Db Schema](https://dbdiagram.io/d/chidi-salon-6780316c6b7fa355c3705d78).
2. [Validator](https://github.com/go-playground/validator/blob/master/README.md).
3. [Go migrate](https://github.com/golang-migrate/migrate/blob/master/README.md).
4. [PG](https://neon.tech/postgresql/postgresql-tutorial/postgresql-foreign-key).
5. [In-memory caching](https://www.codingexplorations.com/blog/harnessing-in-memory-caching-in-go).

## Development cmds
1. `migrate create -ext sql -dir ./database/migrations/ -seq create_users_table`.
2. `docker run --rm -v $(pwd)/database/migrations:/database/migrations migrate/migrate create -ext sql -dir ./database/migrations/ -seq payment_detail`
3. `go clean -testcache`
