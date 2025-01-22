import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import {
  NavigationEnd,
  Router,
  RouterLink,
  RouterOutlet
} from '@angular/router';
import { Breadcrumb } from 'primeng/breadcrumb';
import { MenuItem } from 'primeng/api';
import { RESERVATION_ROUTE } from '@store/store.routes';
import {
  CONFIRM_ROUTE,
  DATES_ROUTE,
  SERVICE_TYPES_ROUTE,
  STAFFS_ROUTE
} from '@store/pages/reservation/reservation.routes';
import { NgClass } from '@angular/common';
import { toSignal } from '@angular/core/rxjs-interop';
import { distinctUntilChanged, filter, map } from 'rxjs';
import { CORE_ROUTE } from '@root/app.routes';
import { STORE_ROUTE } from '@pages/core/core.routes';

@Component({
  selector: 'app-reservation',
  imports: [RouterOutlet, Breadcrumb, RouterLink, NgClass],
  template: ` <div class="w-full">
    <div class="mt-3">
      <p-breadcrumb [model]="items" [home]="home">
        <ng-template #item let-item>
          <a [routerLink]="item.route" class="p-breadcrumb-item-link">
            <span
              class="text-primary font-semibold"
              [ngClass]="{ 'text-primary': endsWith(url(), item.route) }"
              >{{ item.label }}</span
            >
          </a>
        </ng-template>
      </p-breadcrumb>
    </div>

    <div class="pb-3 xl:pb-0 px-3 xl:px-0">
      <router-outlet />
    </div>
  </div>`,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ReservationComponent {
  private readonly router = inject(Router);

  protected readonly url = toSignal(
    this.router.events.pipe(
      filter(event => event instanceof NavigationEnd),
      map(event => event.url),
      distinctUntilChanged()
    ),
    {
      initialValue: `/${CORE_ROUTE}/${STORE_ROUTE}/${RESERVATION_ROUTE}/${SERVICE_TYPES_ROUTE}`
    }
  );

  protected readonly endsWith = (obs: string, crumb: string) => {
    const o = obs.split('/');
    const c = crumb.split('/');
    return o[o.length - 1] === c[c.length - 1];
  };

  protected readonly items: MenuItem[] = [
    {
      label: 'services',
      route: `/${CORE_ROUTE}/${STORE_ROUTE}/${RESERVATION_ROUTE}/${SERVICE_TYPES_ROUTE}`
    },
    {
      label: 'staffs',
      route: `/${CORE_ROUTE}/${STORE_ROUTE}/${RESERVATION_ROUTE}/${STAFFS_ROUTE}`
    },
    {
      label: 'dates',
      route: `/${CORE_ROUTE}/${STORE_ROUTE}/${RESERVATION_ROUTE}/${DATES_ROUTE}`
    },
    {
      label: 'confirm',
      route: `/${CORE_ROUTE}/${STORE_ROUTE}/${RESERVATION_ROUTE}/${CONFIRM_ROUTE}`
    }
  ];

  protected readonly home: MenuItem = {
    label: 'book',
    icon: 'pi pi-shopping-bag',
    routerLink: `/${CORE_ROUTE}/${STORE_ROUTE}/${RESERVATION_ROUTE}/${SERVICE_TYPES_ROUTE}`
  };
}
