import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { ErrorBannerService } from '../services/error-banner.service';
import { catchError, throwError } from 'rxjs';

export const errorInterceptor: HttpInterceptorFn = (req, next) => {
  const router = inject(Router);
  const auth = inject(AuthService);
  const banner = inject(ErrorBannerService);

  return next(req).pipe(
    catchError(err => {
      if (err.status === 401) {
        auth.logout();
        router.navigate(['/login']);
      } else if (err.status === 0) {
        banner.show('Impossibile contattare il server. Controlla la connessione.');
      } else if (err.status >= 500) {
        banner.show('Errore del server, riprova più tardi.');
      }
      return throwError(() => err);
    })
  );
};