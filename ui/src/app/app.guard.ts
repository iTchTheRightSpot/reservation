import { Router } from '@angular/router';
import { inject } from '@angular/core';
import { AuthService } from '@shared/data-access/auth.service';
import { CORE_ROUTE } from './app.routes';
import { LOGIN_ROUTE } from '@pages/core/core.routes';

export const crmGuard = async () => {
  if (!inject(AuthService).activeUser()) {
    await inject(Router).navigate([`${CORE_ROUTE}/${LOGIN_ROUTE}`]);
    return false;
  }
  return true;
};
