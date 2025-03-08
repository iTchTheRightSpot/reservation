import { Routes } from '@angular/router';

export const StaffRoutes = {
  STAFFS: '',
  REGISTER: 'register'
};

export const routes: Routes = [
  {
    path: StaffRoutes.STAFFS,
    loadComponent: () =>
      import('./pages/all/crm-staff.component').then(m => m.CrmStaffComponent)
  },
  {
    path: StaffRoutes.REGISTER,
    loadComponent: () =>
      import('./pages/register/register.component').then(
        m => m.RegisterComponent
      )
  }
];
