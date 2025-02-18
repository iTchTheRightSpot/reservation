import {
  ChangeDetectionStrategy,
  Component,
  inject,
  OnInit,
  signal
} from '@angular/core';
import { ToggleSwitch } from 'primeng/toggleswitch';
import { FormsModule } from '@angular/forms';
import { LightDarkModeService } from '@shared/data-access/light-dark-mode.service';

@Component({
  selector: 'app-light-dark-mode',
  imports: [ToggleSwitch, FormsModule],
  template: `
    <div class="flex gap-1">
      <p-toggleswitch
        [ngModel]="isDarkMode()"
        (ngModelChange)="toggleDarkMode()"
      />
      @if (isDarkMode()) {
        🌚
      } @else {
        🌞
      }
    </div>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class LightDarkModeComponent implements OnInit {
  private readonly service = inject(LightDarkModeService);

  protected readonly isDarkMode = signal(false);

  ngOnInit(): void {
    const el = document.querySelector('html');
    if (
      localStorage.getItem('color-theme') === 'dark' ||
      (!('color-theme' in localStorage) &&
        window.matchMedia('(prefers-color-scheme: dark)').matches)
    ) {
      this.isDarkMode.set(true);
      el?.classList.add('dark');
    } else {
      this.isDarkMode.set(false);
      el?.classList.remove('dark');
    }
    this.service.isDarkMode.set(this.isDarkMode());
  }

  protected readonly toggleDarkMode = () => {
    document.querySelector('html')?.classList.toggle('dark');
    if (this.isDarkMode()) localStorage.setItem('color-theme', 'light');
    else localStorage.setItem('color-theme', 'dark');
    this.isDarkMode.set(!this.isDarkMode());
    this.service.isDarkMode.set(this.isDarkMode());
  };
}
