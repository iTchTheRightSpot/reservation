import { Routes } from '@angular/router';

export const REGISTER_ROUTE = 'register';
export const SETTINGS_ROUTE = 'settings';

export const routes: Routes = [
  {
    path: REGISTER_ROUTE,
    loadComponent: () =>
      import('./pages/register/register.component').then(
        m => m.RegisterComponent
      )
  },
  {
    path: '',
    redirectTo: REGISTER_ROUTE,
    pathMatch: 'full'
  }
];
