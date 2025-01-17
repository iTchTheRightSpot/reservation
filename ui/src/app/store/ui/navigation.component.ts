import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { Menubar } from 'primeng/menubar';
import { MenuItem } from 'primeng/api';
import { Router } from '@angular/router';
import { LightDarkModeComponent } from '@shared/ui/light-dark-mode.component';
import {
  STORE_FRONT_ABOUT_ROUTE,
  STORE_FRONT_RESERVATION_ROUTE
} from '@store/store.routes';
import { STORE_FRONT_HOME_ROUTE } from '@root/app.routes';

@Component({
  selector: 'app-navigation',
  imports: [Menubar, LightDarkModeComponent],
  template: `
    <div class="w-full">
      <p-menubar [model]="items">
        <ng-template #end>
          <app-light-dark-mode />
        </ng-template>
      </p-menubar>
    </div>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class NavigationComponent {
  protected readonly router = inject(Router);

  // TODO comeback to understand why this isn't highlighting active tab
  protected readonly items: MenuItem[] = [
    {
      label: 'HOME',
      icon: 'pi pi-home',
      routerLink: `/${STORE_FRONT_HOME_ROUTE}`,
      routerLinkActiveOptions: { exact: true },
      styleClass:
        this.router.url === `/${STORE_FRONT_HOME_ROUTE}` ? 'text-primary' : ''
    },
    {
      label: 'BOOK',
      icon: 'pi pi-shopping-bag',
      routerLink: `/${STORE_FRONT_RESERVATION_ROUTE}`,
      routerLinkActiveOptions: { exact: true },
      styleClass:
        this.router.url === `/${STORE_FRONT_RESERVATION_ROUTE}`
          ? 'text-primary'
          : ''
    },
    {
      label: 'ABOUT',
      icon: 'pi pi-star',
      routerLink: `/${STORE_FRONT_ABOUT_ROUTE}`,
      routerLinkActiveOptions: { exact: true },
      styleClass:
        this.router.url === `/${STORE_FRONT_ABOUT_ROUTE}` ? 'text-primary' : ''
    }
  ];
}
