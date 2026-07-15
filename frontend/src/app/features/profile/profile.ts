import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { AuthService } from '../../core/services/auth.service';
import { Navbar } from '../../shared/navbar/navbar';

@Component({
  selector: 'app-profile',
  standalone: true,
  imports: [CommonModule, Navbar, FormsModule],
  templateUrl: './profile.html',
  styleUrl: './profile.scss'
})
export class Profile {
  avatarUrl = '';
  saving = false;
  message = '';

  constructor(public auth: AuthService) {
    this.avatarUrl = this.auth.currentUser()?.avatar_url || '';
  }

  saveAvatar() {
    const url = this.avatarUrl.trim();
    this.message = '';

    if (!url) {
      this.message = 'Inserisci un URL';
      return;
    }
    if (!/^https?:\/\//i.test(url)) {
      this.message = "L'URL deve iniziare con http:// o https://";
      return;
    }

    this.saving = true;
    this.validateImage(url).then(isValid => {
      if (!isValid) {
        this.saving = false;
        this.message = "Questo link non punta a un'immagine valida (assicurati di copiare il link diretto al file, non a una pagina web)";
        return;
      }

      this.auth.updateAvatar(url).subscribe({
        next: () => {
          this.saving = false;
          this.message = 'Avatar aggiornato ✅';
        },
        error: () => {
          this.saving = false;
          this.message = 'Errore durante il salvataggio';
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
}