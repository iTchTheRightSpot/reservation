import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { NavigationEnd, Router, RouterOutlet } from '@angular/router';
import { NavigationComponent } from '@store/ui/navigation.component';
import { toSignal } from '@angular/core/rxjs-interop';
import { filter, map, startWith } from 'rxjs';
import { NgClass } from '@angular/common';
import { AuthService } from '@shared/data-access/auth.service';
import { LoadingService } from '@shared/data-access/loading.service';
import { ProgressBar } from 'primeng/progressbar';

@Component({
  selector: 'app-core',
  imports: [RouterOutlet, NavigationComponent, NgClass, ProgressBar],
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
      class="w-full xl:w-cx-50 m-auto z-20"
      [ngClass]="{ pos: url(), pos1: !url() }"
    >
      @if (loading.state()) {
        <p-progressbar mode="indeterminate" [style]="{ height: '6px' }" />
      }
      <app-navigation [imageKey]="auth.activeUser()?.image_key" />
    </div>
    <router-outlet />
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class CoreComponent {
  protected readonly router = inject(Router);
  protected readonly auth = inject(AuthService);
  protected readonly loading = inject(LoadingService);

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
