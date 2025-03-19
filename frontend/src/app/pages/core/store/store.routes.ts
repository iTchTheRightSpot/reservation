import { Routes } from '@angular/router';

export const StoreRoutes = {
  HOME: '',
  RESERVATION: 'reservation',
  ABOUT: 'about'
};

export const routes: Routes = [
  {
    path: StoreRoutes.HOME,
    loadComponent: () =>
      import('./pages/home/home.component').then(m => m.HomeComponent)
  },
  {
    path: StoreRoutes.RESERVATION,
    loadComponent: () =>
      import('./pages/reservation/reservation.component').then(
        m => m.ReservationComponent
      ),
    loadChildren: () =>
      import('./pages/reservation/reservation.routes').then(m => m.routes)
  },
  {
    path: StoreRoutes.ABOUT,
    loadComponent: () =>
      import('./pages/about/about.component').then(m => m.AboutComponent)
  }
];
