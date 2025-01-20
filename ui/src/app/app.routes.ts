import { Routes } from '@angular/router';

export const CORE_ROUTE = '';
export const STAFF_ROUTE = 'staff';
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
    path: STAFF_ROUTE,
    loadComponent: () =>
      import('@root/pages/staff/staff.component').then(m => m.StaffComponent),
    loadChildren: () =>
      import('@root/pages/staff/staff.routes').then(m => m.routes)
  },
  {
    path: NOTFOUND,
    loadComponent: () =>
      import('@root/pages/not-found.component').then(m => m.NotFoundComponent)
  },
  { path: '**', redirectTo: `/${NOTFOUND}` }
];
