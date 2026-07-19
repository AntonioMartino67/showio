import { Injectable, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { environment } from '../../../environments/environment';
import { AuthResponse, User } from '../models/models';
import { tap } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private tokenKey = 'showio_token';
  currentUser = signal<User | null>(null);

  constructor(private http: HttpClient, private router: Router) {}

  register(username: string, email: string, password: string) {
    return this.http.post(`${environment.apiUrl}/register`, { username, email, password });
  }

  verifyOtp(email: string, code: string) {
    return this.http.post<AuthResponse>(`${environment.apiUrl}/verify-otp`, { email, code }).pipe(
      tap(res => {
        localStorage.setItem(this.tokenKey, res.token);
        this.loadUser().subscribe();
      })
    );
  }

  resendOtp(email: string) {
    return this.http.post(`${environment.apiUrl}/resend-otp`, { email });
  }

  login(email: string, password: string) {
    return this.http.post<AuthResponse>(`${environment.apiUrl}/login`, { email, password }).pipe(
      tap(res => {
        localStorage.setItem(this.tokenKey, res.token);
        this.loadUser().subscribe();
      })
    );
  }

  loginWithGoogle(credential: string) {
    return this.http.post<AuthResponse>(`${environment.apiUrl}/auth/google`, { credential }).pipe(
      tap(res => {
        localStorage.setItem(this.tokenKey, res.token);
        this.loadUser().subscribe();
      })
    );
  }

  loadUser() {
    return this.http.get<User>(`${environment.apiUrl}/me`).pipe(
      tap(user => this.currentUser.set(user))
    );
  }

  updateAvatar(avatarUrl: string) {
    return this.http.put(`${environment.apiUrl}/me/avatar`, { avatar_url: avatarUrl }).pipe(
      tap(() => this.currentUser.update(u => u ? { ...u, avatar_url: avatarUrl } : u))
    );
  }

  changePassword(currentPassword: string, newPassword: string) {
    return this.http.put(`${environment.apiUrl}/me/password`, {
      current_password: currentPassword,
      new_password: newPassword
    });
  }

  linkGoogle(credential: string) {
    return this.http.post(`${environment.apiUrl}/me/google`, { credential }).pipe(
      tap(() => this.currentUser.update(u => u ? { ...u, google_linked: true } : u))
    );
  }

  unlinkGoogle() {
    return this.http.delete(`${environment.apiUrl}/me/google`).pipe(
      tap(() => this.currentUser.update(u => u ? { ...u, google_linked: false } : u))
    );
  }

  logout() {
    localStorage.removeItem(this.tokenKey);
    this.currentUser.set(null);
    this.router.navigate(['/login']);
  }

  getToken(): string | null {
    return localStorage.getItem(this.tokenKey);
  }

  isLoggedIn(): boolean {
    return !!this.getToken();
  }

  updateNotifications(notifyNewSeasons: boolean) {
  return this.http.put(`${environment.apiUrl}/me/notifications`, { notify_new_seasons: notifyNewSeasons }).pipe(
    tap(() => this.currentUser.update(u => u ? { ...u, notify_new_seasons: notifyNewSeasons } : u))
  );
}
}