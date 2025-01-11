import { Routes } from '@angular/router';

export const STORE_FRONT_HOME_ROUTE = '';
export const STAFF_HOME_ROUTE = 'staff';
export const NOTFOUND = '404';

export const routes: Routes = [
  {
    path: STORE_FRONT_HOME_ROUTE,
    loadComponent: () =>
      import('./store-front/store-front.component').then(
        m => m.StoreFrontComponent
      ),
    loadChildren: () =>
      import('./store-front/store-front.routes').then(m => m.routes)
  },
  {
    path: STAFF_HOME_ROUTE,
    loadComponent: () =>
      import('./staff/staff.component').then(m => m.StaffComponent),
    loadChildren: () => import('./staff/staff.routes').then(m => m.routes)
  },
  {
    path: NOTFOUND,
    loadComponent: () =>
      import('./shared/pages/not-found.component').then(
        m => m.NotFoundComponent
      )
  },
  { path: '**', redirectTo: `/${NOTFOUND}` }
];
