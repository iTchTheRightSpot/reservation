import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { NavigationEnd, Router, RouterOutlet } from '@angular/router';
import { NavigationComponent } from '@store/ui/navigation.component';
import { toSignal } from '@angular/core/rxjs-interop';
import { filter, map, startWith } from 'rxjs';
import { AuthService } from '@shared/data-access/auth.service';
import { LoadingService } from '@shared/data-access/loading.service';
import { ProgressBar } from 'primeng/progressbar';

@Component({
  selector: 'app-core',
  imports: [RouterOutlet, NavigationComponent, ProgressBar],
  template: `
    <div class="w-full xl:max-w-7xl m-auto">
      <div class="w-full sticky top-0 z-20">
        @if (loading.state()) {
          <p-progress-bar mode="indeterminate" [style]="{ height: '6px' }" />
        }
        <app-navigation [imageKey]="auth.activeUser()?.image_key" />
      </div>

      <div class="w-full pb-2 px-1">
        <router-outlet />
      </div>
    </div>
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
