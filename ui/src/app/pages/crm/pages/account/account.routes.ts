import { Routes } from '@angular/router';

export const REGISTER_ROUTE = 'register';
export const SETTINGS_ROUTE = 'settings';
export const SCHEDULE_ROUTE = 'schedule';

export const routes: Routes = [
  {
    path: REGISTER_ROUTE,
    loadComponent: () =>
      import('./pages/register/register.component').then(
        m => m.RegisterComponent
      )
  },
  {
    path: SCHEDULE_ROUTE,
    loadComponent: () =>
      import('./pages/schedule/schedule.component').then(
        m => m.ScheduleComponent
      )
  },
  {
    path: '',
    redirectTo: REGISTER_ROUTE,
    pathMatch: 'full'
  }
];
