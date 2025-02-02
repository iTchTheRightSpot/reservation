import { Routes } from '@angular/router';
import { crmGuard } from './app.guard';

export const CORE_ROUTE = '';
export const CRM_ROUTE = 'crm';
export const NOTFOUND = '404';

export const routes: Routes = [
  {
    path: CORE_ROUTE,
    loadComponent: () =>
      import('@root/pages/core/core.component').then(m => m.CoreComponent),
    loadChildren: () =>
      import('@root/pages/core/core.routes').then(m => m.routes)
  },
  {
    path: CRM_ROUTE,
    canActivate: [crmGuard],
    canActivateChild: [crmGuard],
    loadComponent: () =>
      import('@pages/crm/crm.component').then(m => m.CrmComponent),
    loadChildren: () => import('@pages/crm/crm.routes').then(m => m.routes)
  },
  {
    path: NOTFOUND,
    loadComponent: () =>
      import('@root/pages/not-found.component').then(m => m.NotFoundComponent)
  },
  { path: '**', redirectTo: `/${NOTFOUND}` }
];
