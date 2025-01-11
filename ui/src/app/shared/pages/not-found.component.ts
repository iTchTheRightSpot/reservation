import { ChangeDetectionStrategy, Component } from '@angular/core';
import { STORE_FRONT_HOME_ROUTE } from '@root/app.routes';
import { RouterLink } from '@angular/router';
import { Button } from 'primeng/button';

@Component({
  selector: 'app-not-found',
  imports: [RouterLink, Button],
  template: `
    <section class="min-w-screen min-h-screen flex items-center justify-center">
      <div
        class="p-5 sm:p-20 flex gap-8 flex-col items-center rounded-md border border-solid border-transparent bg-[#18181b]"
      >
        <h1 class="text-7xl tracking-tight font-extrabold lg:text-9xl">404</h1>
        <div
          class="font-bold text-center text-4xl border-t border-surface pt-8"
        >
          Page Not Found
        </div>
        <p-button
          label="GO TO HOMEPAGE"
          [routerLink]="HOME"
          tabindex="0"
        ></p-button>
      </div>
    </section>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class NotFoundComponent {
  protected readonly HOME = STORE_FRONT_HOME_ROUTE;
}
