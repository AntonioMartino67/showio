import { Component, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router, ActivatedRoute, RouterLink } from '@angular/router';
import { AuthService } from '../../../core/services/auth.service';

@Component({
  selector: 'app-verify-otp',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink],
  templateUrl: './verify-otp.html',
  styleUrl: './verify-otp.scss'
})
export class VerifyOtp {
  email = '';
  code = '';
  error = signal('');
  info = signal('');
  loading = signal(false);
  resending = signal(false);

  constructor(private auth: AuthService, private router: Router, private route: ActivatedRoute) {
    this.email = this.route.snapshot.queryParamMap.get('email') || '';
  }

  submit() {
    this.error.set('');
    this.info.set('');
    this.loading.set(true);
    this.auth.verifyOtp(this.email, this.code).subscribe({
      next: () => {
        this.loading.set(false);
        this.router.navigate(['/dashboard']);
      },
      error: (err) => {
        this.loading.set(false);
        this.error.set(err.error || 'Codice non valido');
      }
    });
  }

  resend() {
    this.error.set('');
    this.info.set('');
    this.resending.set(true);
    this.auth.resendOtp(this.email).subscribe({
      next: () => {
        this.resending.set(false);
        this.info.set('Nuovo codice inviato');
      },
      error: (err) => {
        this.resending.set(false);
        this.error.set(err.error || 'Invio fallito');
      }
    });
  }
}