# Reservation application powered by Go & PostgresSQL

## About
Service similar to Square reservation

## Core Requirements

1. **Multi-tenant**
   - Single database, multiple schema

2. **Personal Identifiable Information**
   - A user should be able to delete their personal info but only after 1 year of inactivity.
   - Data integrity should not be affected.

3. **Roles and Permissions**
    - Data should be protected and only allowed to users with specific roles and or permissions.

4. **Staff**
    - Bares minimum to be a staff is to have a role `STAFF`.
    - **Schedule**
      - Only staffs with `WRITE` permission can create working hours for other staffs.
      - Allow bulk creation of schedules with the following constraints:
      - No conflicts in working hrs.
      - A schedule can only be deleted if not reservation is attached to it, but it can be as not
      visible.
      - Schedule cannot be created for past date.
      - **Autonomous Weekly Recurring Schedule**
        - Similar to marking a reoccurring alarm, autonomous schedule that runs weekly, for schedules
        that have been marked as reoccurring. Cron would fire once a week.
   - **Services**
      - Only staffs with `WRITE or DELETE` permission can create or delete a service.
      - A service cannot be deleted if it has an existing relationship with another table.
      Instead, it's visibility can be marked as false so clients can reserve for said service.
      - Only staff with `WRITE` can assign a service to another staff.
    - **Reservations**
      - Only staffs with `WRITE` permission can cancel a reservation.
      - A reservation can include multiple services.
      - A reservation can can have an amount quoted.
      - No overbooking or time overlap for n number of services for a specific staff.
      - A reservation can never be deleted rather cancelled.
      - Said staff should receive a notification on reservation creation and status change.
      - For a reservation status to be marked as `COMPLETED`, it has to have a payment detail.
      That is an existing relationship with payment_detail.
      - Both staff & client should receive a notification 1 day before a `PENDING` appointment.

5. **Clients**
   - **Services**
     - A client should see all the services offered.
     - n number of services can be for one reservation.
     - Every service can have a max price of `DECIMAL(6, 2)`.
   - **Reservations**
     - A client can only be shown n valid reservation times for 1 staff.
     - Although business can be in a different timezone, valid reservation times should match out to
     said clients timezone.
     - Clients cannot reschedule a reservation. The reservation has to be cancelled.
     - A reservation cannot be made for past dates.
     - Clients & staffs should receive notifications on appointment status change.

6. **Payment**
    - Only staffs can send invoices to clients.
    - Clients should be able to pay via online invoice.

# Development docs
1. [Db Schema](https://dbdiagram.io/d/landscape-erp-66303ee65b24a634d01e83ea).
2. [Validator](https://github.com/go-playground/validator/blob/master/README.md).
3. [Go migrate](https://github.com/golang-migrate/migrate/blob/master/README.md).
4. [PG](https://neon.tech/postgresql/postgresql-tutorial/postgresql-foreign-key).
5. [In-memory caching](https://www.codingexplorations.com/blog/harnessing-in-memory-caching-in-go).

## Cmd
1. `migrate create -ext sql -dir ./database/migrations/ -seq create_users_table`.
2. `go clean -testcache`