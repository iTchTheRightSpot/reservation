import { Routes } from '@angular/router';

export const STORE_ROUTE = '';
export const LOGIN_ROUTE = 'login';

export const routes: Routes = [
  {
    path: STORE_ROUTE,
    loadComponent: () =>
      import('@pages/core/store/store.component').then(m => m.StoreComponent),
    loadChildren: () =>
      import('@pages/core/store/store.routes').then(m => m.routes)
  },
  {
    path: LOGIN_ROUTE,
    loadComponent: () =>
      import('@pages/core/login/login.component').then(m => m.LoginComponent)
  }
];
