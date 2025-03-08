import { Router, Routes } from '@angular/router';
import { inject } from '@angular/core';
import { AuthService } from '@shared/data-access/auth.service';
import { CoreRoutes } from '@core/core.routes';

export const RootRoutes = {
  CORE: '',
  CRM: 'crm',
  NOTFOUND: '404'
};

export const crmGuard = async () => {
  if (!inject(AuthService).activeUser()) {
    await inject(Router).navigate([`${RootRoutes.CORE}/${CoreRoutes.LOGIN}`]);
    return false;
  }
  return true;
};

export const routes: Routes = [
  {
    path: RootRoutes.CORE,
    loadComponent: () =>
      import('@core/core.component').then(m => m.CoreComponent),
    loadChildren: () => import('@core/core.routes').then(m => m.routes)
  },
  {
    path: RootRoutes.CRM,
    canActivate: [crmGuard],
    canActivateChild: [crmGuard],
    loadComponent: () => import('@crm/crm.component').then(m => m.CrmComponent),
    loadChildren: () => import('@crm/crm.routes').then(m => m.routes)
  },
  {
    path: RootRoutes.NOTFOUND,
    loadComponent: () =>
      import('@pages/not-found.component').then(m => m.NotFoundComponent)
  },
  { path: '**', redirectTo: `/${RootRoutes.NOTFOUND}` }
];
