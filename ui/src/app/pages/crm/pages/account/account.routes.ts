import { Routes } from '@angular/router';

export const SCHEDULE_ROUTE = 'schedule';
export const SETTINGS_ROUTE = 'settings';

export const routes: Routes = [
  {
    path: SCHEDULE_ROUTE,
    loadComponent: () =>
      import('./pages/schedule/schedule.component').then(
        m => m.ScheduleComponent
      )
  },
  {
    path: SETTINGS_ROUTE,
    loadComponent: () =>
      import('./pages/settings/settings.component').then(
        m => m.SettingsComponent
      )
  },
  {
    path: '',
    redirectTo: SCHEDULE_ROUTE,
    pathMatch: 'full'
  }
];
