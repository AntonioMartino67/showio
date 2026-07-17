import { Component, signal, AfterViewInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { AuthService } from '../../core/services/auth.service';
import { Navbar } from '../../shared/navbar/navbar';
import { PasswordInput } from '../../shared/password-input/password-input';
import { environment } from '../../../environments/environment';

declare const google: any;

@Component({
  selector: 'app-profile',
  standalone: true,
  imports: [CommonModule, Navbar, FormsModule, PasswordInput],
  templateUrl: './profile.html',
  styleUrl: './profile.scss'
})
export class Profile implements AfterViewInit {
  
  avatarUrl = '';
  saving = signal(false);
  message = signal('');
  isError = signal(false);

  constructor(public auth: AuthService) {
    this.avatarUrl = this.auth.currentUser()?.avatar_url || '';
  }

  saveAvatar() {
    const url = this.avatarUrl.trim();
    this.message.set('');
    this.isError.set(false);

    if (!url) {
      this.message.set('Inserisci un URL');
      this.isError.set(true);
      return;
    }
    if (!/^https?:\/\//i.test(url)) {
      this.message.set("L'URL deve iniziare con http:// o https://");
      this.isError.set(true);
      return;
    }

    this.saving.set(true);
    this.validateImage(url).then(isValid => {
      if (!isValid) {
        this.saving.set(false);
        this.isError.set(true);
        this.message.set("Questo link non punta a un'immagine valida (assicurati di copiare il link diretto al file, non a una pagina web)");
        return;
      }

      this.auth.updateAvatar(url).subscribe({
        next: () => {
          this.saving.set(false);
          this.isError.set(false);
          this.message.set('Avatar aggiornato ✅');
        },
        error: () => {
          this.saving.set(false);
          this.isError.set(true);
          this.message.set('Errore durante il salvataggio');
        }
      });
    });
  }

  private validateImage(url: string): Promise<boolean> {
    return new Promise(resolve => {
      const img = new Image();
      const timeout = setTimeout(() => resolve(false), 6000);
      img.onload = () => { clearTimeout(timeout); resolve(true); };
      img.onerror = () => { clearTimeout(timeout); resolve(false); };
      img.src = url;
    });
  }

  currentPassword = '';
  newPassword = '';
  confirmPassword = '';
  pwSaving = signal(false);
  pwMessage = signal('');
  pwIsError = signal(false);

  changePassword() {
    this.pwMessage.set('');
    this.pwIsError.set(false);

    const hasPassword = this.auth.currentUser()?.has_password;

    if (hasPassword && !this.currentPassword) {
      this.pwIsError.set(true);
      this.pwMessage.set('Inserisci la password attuale');
      return;
    }
    if (this.newPassword.length < 6) {
      this.pwIsError.set(true);
      this.pwMessage.set('La nuova password deve avere almeno 6 caratteri');
      return;
    }
    if (this.newPassword !== this.confirmPassword) {
      this.pwIsError.set(true);
      this.pwMessage.set('Le password non coincidono');
      return;
    }

    this.pwSaving.set(true);
    this.auth.changePassword(this.currentPassword, this.newPassword).subscribe({
      next: () => {
        this.pwSaving.set(false);
        this.pwIsError.set(false);
        this.pwMessage.set('Password aggiornata ✅');
        this.currentPassword = '';
        this.newPassword = '';
        this.confirmPassword = '';
        this.auth.currentUser.update(u => u ? { ...u, has_password: true } : u);
      },
      error: (err) => {
        this.pwSaving.set(false);
        this.pwIsError.set(true);
        if (err?.status === 403) {
          this.pwMessage.set('Password attuale errata');
        } else if (err?.status === 409) {
          this.pwMessage.set('Non puoi riutilizzare una delle ultime 5 password');
        } else {
          this.pwMessage.set('Errore durante il salvataggio');
        }
      }
    });
  }

  googleMessage = signal('');
  googleIsError = signal(false);

  ngAfterViewInit() {
    if (this.auth.currentUser()?.google_linked) return;
    if (typeof google === 'undefined') return;
    google.accounts.id.initialize({
      client_id: environment.googleClientId,
      callback: (response: any) => this.handleGoogleLink(response)
    });
    google.accounts.id.renderButton(
      document.getElementById('google-link-btn'),
      { theme: 'outline', size: 'large', width: 280 }
    );
  }

  handleGoogleLink(response: any) {
    this.googleMessage.set('');
    this.auth.linkGoogle(response.credential).subscribe({
      next: () => {
        this.googleIsError.set(false);
        this.googleMessage.set('Account Google collegato ✅');
      },
      error: (err) => {
        this.googleIsError.set(true);
        this.googleMessage.set(err?.status === 409 ? 'Questo account Google è già collegato a un altro utente' : 'Collegamento fallito');
      }
    });
  }

  unlinkGoogle() {
    this.googleMessage.set('');
    this.auth.unlinkGoogle().subscribe({
      next: () => {
        this.googleIsError.set(false);
        this.googleMessage.set('Account Google scollegato');
      },
      error: (err) => {
        this.googleIsError.set(true);
        this.googleMessage.set(err?.status === 409 ? 'Imposta prima una password, altrimenti perderesti l\'accesso' : 'Scollegamento fallito');
      }
    });
  }
}