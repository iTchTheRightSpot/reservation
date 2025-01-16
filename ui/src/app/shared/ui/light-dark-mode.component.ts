import { ChangeDetectionStrategy, Component, signal } from '@angular/core';
import { ToggleSwitch } from 'primeng/toggleswitch';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-light-dark-mode',
  imports: [ToggleSwitch, FormsModule],
  template: `
    <div class="flex gap-1">
      @if (isDarkMode()) {
        🌚
      } @else {
        🌞
      }
      <p-toggleswitch
        [ngModel]="isDarkMode()"
        (ngModelChange)="toggleDarkMode()"
      />
    </div>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class LightDarkModeComponent {
  protected readonly isDarkMode = signal(false);

  protected readonly toggleDarkMode = () => {
    const element = document.querySelector('html');
    element?.classList.toggle('dark');
    this.isDarkMode.set(!this.isDarkMode());
  };
}
