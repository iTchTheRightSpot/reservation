import { Routes } from '@angular/router';

export const AccountRoutes = {
  SCHEDULE: 'schedule',
  SETTINGS: 'settings'
};

export const routes: Routes = [
  {
    path: AccountRoutes.SCHEDULE,
    loadComponent: () =>
      import('./pages/schedule/schedule.component').then(
        m => m.ScheduleComponent
      )
  },
  {
    path: AccountRoutes.SETTINGS,
    loadComponent: () =>
      import('./pages/settings/settings.component').then(
        m => m.SettingsComponent
      )
  },
  {
    path: '',
    redirectTo: AccountRoutes.SCHEDULE,
    pathMatch: 'full'
  }
];
