import { Routes } from '@angular/router';
import { STORE_FRONT_HOME_ROUTE } from '@root/app.routes';

export const STORE_FRONT_RESERVATION_ROUTE = 'reservation';
export const STORE_FRONT_ABOUT_ROUTE = 'about';

export const routes: Routes = [
  {
    path: STORE_FRONT_HOME_ROUTE,
    loadComponent: () =>
      import('@store/pages/home/home.component').then(m => m.HomeComponent)
  },
  {
    path: STORE_FRONT_RESERVATION_ROUTE,
    loadComponent: () =>
      import('@store/pages/reservation/reservation.component').then(
        m => m.ReservationComponent
      ),
    loadChildren: () =>
      import('@store/pages/reservation/reservation.routes').then(m => m.routes)
  },
  {
    path: STORE_FRONT_ABOUT_ROUTE,
    loadComponent: () =>
      import('@store/pages/about/about.component').then(m => m.AboutComponent)
  }
];
