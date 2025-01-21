import { Routes } from '@angular/router';

export const DASHBOARD_ROUTE = 'dashboard';
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
    path: ACCOUNT_ROUTE,
    loadComponent: () =>
      import('@crm/pages/account/account.component').then(
        m => m.AccountComponent
      ),
    loadChildren: () =>
      import('@crm/pages/account/account.routes').then(m => m.routes)
  },
  {
    path: '',
    redirectTo: DASHBOARD_ROUTE,
    pathMatch: 'full'
  }
];
