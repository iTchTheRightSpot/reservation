import { Routes } from '@angular/router';

export const HOME_ROUTE = '';
export const RESERVATION_ROUTE = 'reservation';
export const ABOUT_ROUTE = 'about';

export const routes: Routes = [
  {
    path: HOME_ROUTE,
    loadComponent: () =>
      import('@store/pages/home/home.component').then(m => m.HomeComponent)
  },
  {
    path: RESERVATION_ROUTE,
    loadComponent: () =>
      import('@store/pages/reservation/reservation.component').then(
        m => m.ReservationComponent
      ),
    loadChildren: () =>
      import('@store/pages/reservation/reservation.routes').then(m => m.routes)
  },
  {
    path: ABOUT_ROUTE,
    loadComponent: () =>
      import('@store/pages/about/about.component').then(m => m.AboutComponent)
  }
];
