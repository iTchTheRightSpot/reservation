import { ChangeDetectionStrategy, Component, input } from '@angular/core';
import { Menubar } from 'primeng/menubar';
import { MenuItem } from 'primeng/api';
import { LightDarkModeComponent } from '@shared/ui/light-dark-mode.component';
import { ABOUT_ROUTE, RESERVATION_ROUTE } from '@store/store.routes';
import { CORE_ROUTE, CRM_ROUTE } from '@root/app.routes';
import { Avatar } from 'primeng/avatar';
import { LOGIN_ROUTE, STORE_ROUTE } from '@pages/core/core.routes';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-navigation',
  imports: [Menubar, LightDarkModeComponent, Avatar, RouterLink],
  template: `
    <div class="w-full">
      <p-menubar [model]="items">
        <ng-template #end>
          <div class="flex gap-2 items-center justify-center">
            @if (imageKey(); as k) {
              <a [routerLink]="crm">
                <p-avatar [image]="k" class="mr-2" shape="circle" />
              </a>
            } @else if (imageKey() === null) {
              <a [routerLink]="crm">
                <p-avatar icon="pi pi-user" class="mr-2" shape="circle" />
              </a>
            }
            <app-light-dark-mode />
          </div>
        </ng-template>
      </p-menubar>
    </div>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class NavigationComponent {
  imageKey = input.required<string | null | undefined>();

  protected readonly crm = `${CRM_ROUTE}`;

  protected readonly items: MenuItem[] = [
    {
      label: 'Home',
      icon: 'pi pi-home',
      routerLink: `/${CORE_ROUTE}/${STORE_ROUTE}/${STORE_ROUTE}`
    },
    {
      label: 'Book',
      icon: 'pi pi-shopping-bag',
      routerLink: `/${CORE_ROUTE}/${STORE_ROUTE}/${RESERVATION_ROUTE}`
    },
    {
      label: 'About',
      icon: 'pi pi-file',
      routerLink: `/${CORE_ROUTE}/${STORE_ROUTE}/${ABOUT_ROUTE}`
    },
    {
      label: 'Login',
      icon: 'pi pi-user',
      routerLink: `/${CORE_ROUTE}/${LOGIN_ROUTE}`
    }
  ];
}
