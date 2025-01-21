import { Routes } from '@angular/router';

export const DASHBOARD_ROUTE = 'dashboard';
export const CRM_SERVICE_TYPE_ROUTE = 'service';
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
    path: CRM_SERVICE_TYPE_ROUTE,
    loadComponent: () =>
      import('./pages/service-type/crm-service-type.component').then(
        m => m.CRMServiceTypeComponent
      )
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
