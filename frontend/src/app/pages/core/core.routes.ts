import { Routes } from '@angular/router';

export const CoreRoutes = {
  STORE: '',
  LOGIN: 'login'
};

export const routes: Routes = [
  {
    path: CoreRoutes.STORE,
    loadComponent: () =>
      import('./store/store.component').then(m => m.StoreComponent),
    loadChildren: () => import('./store/store.routes').then(m => m.routes)
  },
  {
    path: CoreRoutes.LOGIN,
    loadComponent: () =>
      import('./login/login.component').then(m => m.LoginComponent)
  }
];
