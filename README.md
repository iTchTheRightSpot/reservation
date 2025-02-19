# Reservation application powered by Angular, Go & PostgresSQL

## About
Service similar to Square reservation

## Core Requirements

1. **Authentication & Authorization**
   - Authentication by email & password
   - Role & permission based authorization using jwt
     - To access protected resources, each request must contain jwt as a cookie.

2. **Staff**
    - **Schedule**
      - Only staffs with `WRITE` permission can create working hours for other staffs.
      - Allow bulk creation of schedules with the following constraints:
      - No conflicts in working hrs.
      - A schedule can only be deleted if no reservation is attached to it.
      - Schedules cannot be created for past date.
      - **Autonomous Weekly Recurring Schedule**
        - Similar to marking a reoccurring alarm, autonomous schedule that runs weekly, for schedules
        that have been marked as reoccurring. Cron would fire once a week.
   - **Services**
      - Only staffs with `WRITE or DELETE` permission can create or delete a service.
      - A service cannot be deleted if it has an existing relationship with another table.
      Instead, it's visibility can be marked as false so clients can reserve for said service.
      - Only staff with `WRITE` can assign a service to another staff.
    - **Reservations**
      - A reservation can include multiple services.
      - A reservation must have a quoted amount.
      - No overbooking.
      - A reservation can never be deleted rather cancelled.
      - Said staff should receive a notification on reservation creation and status change.
      - For a reservation status to be marked as `COMPLETED`, it has to have a payment detail.
      That is an existing relationship with payment_detail.
      - Both staff & client should receive a notification 1 day before a `PENDING` appointment.

3. **Clients**
   - **Services**
     - A client should see all the services offered.
     - n number of services can be for one reservation.
     - Every service can have a max price of `DECIMAL(6, 2)`.
   - **Reservations**
     - A client can only be shown n valid reservation times for 1 staff.
     - Reservation times should be adjusted to a clients' timezone.
     - Clients cannot reschedule a reservation. The reservation has to be cancelled.
     - A reservation cannot be made for past dates.
     - Clients & staffs should receive notifications on appointment status change.

4. **Payment**
    - Only staffs can send invoices to clients.
    - Clients should be able to pay via online invoice.

# Development docs
1. [Db Schema](https://dbdiagram.io/d/chidi-salon-6780316c6b7fa355c3705d78).
2. [Validator](https://github.com/go-playground/validator/blob/master/README.md).
3. [Go migrate](https://github.com/golang-migrate/migrate/blob/master/README.md).
4. [PG](https://neon.tech/postgresql/postgresql-tutorial/postgresql-foreign-key).
5. [In-memory caching](https://www.codingexplorations.com/blog/harnessing-in-memory-caching-in-go).

## Development details
1. `migrate create -ext sql -dir ./database/migrations/ -seq create_users_table`.
2. `docker run --rm -v $(pwd)/database/migrations:/database/migrations migrate/migrate create -ext sql -dir ./database/migrations/ -seq payment_detail`
3. `go clean -testcache`