import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { NavigationComponent } from '@store/ui/navigation.component';
import { NavigationEnd, Router, RouterOutlet } from '@angular/router';
import { NgClass } from '@angular/common';
import { filter, map, startWith } from 'rxjs';
import { toSignal } from '@angular/core/rxjs-interop';

@Component({
  selector: 'app-store-front',
  imports: [NavigationComponent, RouterOutlet, NgClass],
  styles: [
    `
      .pos {
        position: fixed;
        left: 0;
        top: 0;
        right: 0;
      }

      .pos1 {
        position: sticky;
        top: 0;
      }
    `
  ],
  template: `
    <div
      class="w-full xl:w-cx-50 m-auto z-10"
      [ngClass]="{ pos: url(), pos1: !url() }"
    >
      <app-navigation />
    </div>
    <router-outlet />
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class StoreComponent {
  protected readonly router = inject(Router);

  protected readonly url = toSignal(
    this.router.events.pipe(
      filter(event => event instanceof NavigationEnd),
      map(event => event.url),
      startWith(this.router.url),
      map(e => e === '/')
    ),
    { initialValue: false }
  );
}
