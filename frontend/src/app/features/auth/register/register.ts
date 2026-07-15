import { Component, AfterViewInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../../../core/services/auth.service';
import { environment } from '../../../../environments/environment';
import { PasswordInput } from '../../../shared/password-input/password-input';

declare const google: any;

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink, PasswordInput],
  templateUrl: './register.html',
  styleUrl: './register.scss'
})
export class Register implements AfterViewInit {
  username = '';
  email = '';
  password = '';
  confirmPassword = '';
  error = signal('');
  loading = signal(false);

  constructor(private auth: AuthService, private router: Router) {}

  ngAfterViewInit() {
    if (typeof google === 'undefined') return;
    google.accounts.id.initialize({
      client_id: environment.googleClientId,
      callback: (response: any) => this.handleGoogleResponse(response)
    });
    google.accounts.id.renderButton(
      document.getElementById('google-btn'),
      { theme: 'outline', size: 'large', width: 320 }
    );
  }

  handleGoogleResponse(response: any) {
    this.error.set('');
    this.loading.set(true);
    this.auth.loginWithGoogle(response.credential).subscribe({
      next: () => {
        this.loading.set(false);
        this.router.navigate(['/dashboard']);
      },
      error: () => {
        this.loading.set(false);
        this.error.set('Registrazione con Google fallita');
      }
    });
  }

  submit() {
    this.error.set('');
    if (this.password !== this.confirmPassword) {
      this.error.set('Le password non coincidono');
      return;
    }
    this.loading.set(true);
    this.auth.register(this.username, this.email, this.password).subscribe({
      next: () => {
        this.loading.set(false);
        this.router.navigate(['/verify-otp'], { queryParams: { email: this.email } });
      },
      error: (err) => {
        this.loading.set(false);
        this.error.set(err.error || 'Errore durante la registrazione');
      }
    });
  }
}