import { Injectable, signal } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class LightDarkModeService {
  readonly isDarkMode = signal(false);
}
