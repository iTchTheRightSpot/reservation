import { Routes } from '@angular/router';

export const ALL_STAFFS_ROUTE = '';
export const REGISTER_ROUTE = 'register';

export const routes: Routes = [
  {
    path: ALL_STAFFS_ROUTE,
    loadComponent: () =>
      import('./pages/all-staffs/all-staffs.component').then(
        m => m.AllStaffsComponent
      )
  },
  {
    path: REGISTER_ROUTE,
    loadComponent: () =>
      import('./pages/register/register.component').then(
        m => m.RegisterComponent
      )
  }
];
