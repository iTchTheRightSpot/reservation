import { Routes } from '@angular/router';

export const DASHBOARD_ROUTE = 'dashboard';
export const CRM_SERVICE_TYPES_ROUTE = 'services';
export const BOOKINGS_ROUTE = 'bookings';
export const CRM_STAFFS_ROUTE = 'staffs';
export const ACCOUNT_ROUTE = 'account';

export const routes: Routes = [
  {
    path: DASHBOARD_ROUTE,
    loadComponent: () =>
      import('./pages/dashboard/dashboard.component').then(
        m => m.DashboardComponent
      )
  },
  {
    path: BOOKINGS_ROUTE,
    loadComponent: () =>
      import('./pages/bookings/bookings.component').then(
        m => m.BookingsComponent
      )
  },
  {
    path: CRM_SERVICE_TYPES_ROUTE,
    loadComponent: () =>
      import('./pages/service-type/crm-service-type.component').then(
        m => m.CRMServiceTypeComponent
      )
  },
  {
    path: CRM_STAFFS_ROUTE,
    loadComponent: () =>
      import('./pages/staff/staff-core.component').then(
        m => m.StaffCoreComponent
      ),
    loadChildren: () =>
      import('./pages/staff/staff-core.routes').then(m => m.routes)
  },
  {
    path: ACCOUNT_ROUTE,
    loadComponent: () =>
      import('./pages/account/account.component').then(m => m.AccountComponent),
    loadChildren: () =>
      import('./pages/account/account.routes').then(m => m.routes)
  },
  {
    path: '',
    redirectTo: DASHBOARD_ROUTE,
    pathMatch: 'full'
  }
];
