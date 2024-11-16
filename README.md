# Landscape ERP application powered by Go

## Core Requirements

1. **Roles and Permissions**
    - Staff members should have different roles.
    - Staff members should have different permissions. For example, not all staffs should have
   `WRITE` permissions.

2. **Working Hour Creation**
    - Only staff members with the `WRITE` permission can create working hours for other staff.

3. **Bulk Shift Creation**
    - Allow bulk creation of shifts with the following constraints:
        - No duplicate shifts.
        - Shifts cannot be created for past dates.

4. **Autonomous Weekly Recurring Schedule**
    - Implement an autonomous schedule that runs weekly, potentially as a cron job every Sunday.

5. **Recurring Shift Flag**
    - Add an `isRecurring` field to shifts. The cron job should:
        - Retrieve shifts from the past 7 days when `isRecurring` is true.
        - Recreate these recurring shifts.

6. **Client Scheduling**
    - Clients should be able to schedule a time based on the working hours of staff.
    - No two clients can reserve the same time for a staff member.

7. **Timezone Handling**
    - Time slots for clients should be reflected in the client's timezone.

8. **Invoice Permissions**
    - Only staff members with the appropriate permissions can send invoices to clients.

9. **Payment Flexibility**
    - Some appointments/jobs require a partial (half) prepayment, while others do not.
    - Although an appointment has a fixed price, payments can be split across different times.

# Development docs
1. [Db Schema](https://dbdiagram.io/d/landscape-erp-66303ee65b24a634d01e83ea).
2. [Validator](https://github.com/go-playground/validator/blob/master/README.md).
3. [Go migrate](https://github.com/golang-migrate/migrate/blob/master/README.md).
4. [PG](https://neon.tech/postgresql/postgresql-tutorial/postgresql-foreign-key).
5. [In-memory caching](https://www.codingexplorations.com/blog/harnessing-in-memory-caching-in-go).

## Cmd
1. `migrate create -ext sql -dir ./database/migrations/ -seq create_users_table`.
