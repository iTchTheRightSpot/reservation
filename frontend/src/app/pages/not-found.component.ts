import { ChangeDetectionStrategy, Component } from '@angular/core';
import { RootRoutes } from '@root/app.routes';
import { RouterLink } from '@angular/router';
import { Button } from 'primeng/button';

@Component({
  selector: 'app-not-found',
  imports: [RouterLink, Button],
  template: `
    <section
      class="bg-[#18181b] min-w-screen min-h-screen flex items-center justify-center"
    >
      <div
        class="p-5 sm:p-20 flex gap-8 flex-col items-center rounded-md border border-solid border-transparent"
      >
        <h1
          class="text-white text-7xl tracking-tight font-extrabold lg:text-9xl"
        >
          404
        </h1>
        <div
          class="font-bold text-center text-4xl border-t border-surface pt-8"
        >
          <p class="text-white">Page Not Found</p>
        </div>
        <p-button label="GO TO HOMEPAGE" [routerLink]="HOME" tabindex="0" />
      </div>
    </section>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class NotFoundComponent {
  protected readonly HOME = RootRoutes.CORE;
}
