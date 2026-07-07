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

  login(email: string, password: string) {
    return this.http.post<AuthResponse>(`${environment.apiUrl}/login`, { email, password }).pipe(
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
}