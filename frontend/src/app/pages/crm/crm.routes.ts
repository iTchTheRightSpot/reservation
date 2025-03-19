import { Routes } from '@angular/router';

export const CRMRoutes = {
  DASHBOARD: 'dashboard',
  SERVICES: 'services',
  BOOKINGS: 'bookings',
  STAFFS: 'staffs',
  ACCOUNT: 'account'
};

export const routes: Routes = [
  {
    path: CRMRoutes.DASHBOARD,
    loadComponent: () =>
      import('./pages/dashboard/dashboard.component').then(
        m => m.DashboardComponent
      )
  },
  {
    path: CRMRoutes.BOOKINGS,
    loadComponent: () =>
      import('./pages/bookings/bookings.component').then(
        m => m.BookingsComponent
      )
  },
  {
    path: CRMRoutes.SERVICES,
    loadComponent: () =>
      import('./pages/service-type/crm-service-type.component').then(
        m => m.CRMServiceTypeComponent
      )
  },
  {
    path: CRMRoutes.STAFFS,
    loadComponent: () =>
      import('./pages/staff/staff-core.component').then(
        m => m.StaffCoreComponent
      ),
    loadChildren: () =>
      import('./pages/staff/staff-core.routes').then(m => m.routes)
  },
  {
    path: CRMRoutes.ACCOUNT,
    loadComponent: () =>
      import('./pages/account/account.component').then(m => m.AccountComponent),
    loadChildren: () =>
      import('./pages/account/account.routes').then(m => m.routes)
  },
  {
    path: '',
    redirectTo: CRMRoutes.DASHBOARD,
    pathMatch: 'full'
  }
];
