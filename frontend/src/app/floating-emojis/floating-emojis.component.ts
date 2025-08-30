import { Component, OnInit, OnDestroy } from '@angular/core';

@Component({
  selector: 'app-floating-emojis',
  template: `
    <div *ngFor="let emoji of emojis; trackBy: trackByIndex" 
         class="floating-emoji" 
         [style.top.px]="emoji.y"
         [style.left.px]="emoji.x"
         [style.animation-delay.s]="emoji.delay">
      {{emoji.symbol}}
    </div>
  `,
  styleUrls: ['./floating-emojis.component.css']
})
export class FloatingEmojisComponent implements OnInit, OnDestroy {
  emojis: Array<{symbol: string, x: number, y: number, delay: number}> = [];
  private intervalId: any;

  private emojiSymbols = ['🕵️', '🐱', '🔫', '🐭', '🗝️', '💎', '🕶️', '🎯', '⚡', '🌟'];

  ngOnInit() {
    this.generateEmojis();
    this.startAnimation();
  }

  ngOnDestroy() {
    if (this.intervalId) {
      clearInterval(this.intervalId);
    }
  }

  private generateEmojis() {
    const numberOfEmojis = Math.min(8, Math.max(4, Math.floor(window.innerWidth / 300)));
    
    for (let i = 0; i < numberOfEmojis; i++) {
      this.emojis.push({
        symbol: this.getRandomEmoji(),
        x: Math.random() * (window.innerWidth - 100),
        y: Math.random() * (window.innerHeight - 100),
        delay: Math.random() * 6
      });
    }
  }

  private getRandomEmoji(): string {
    return this.emojiSymbols[Math.floor(Math.random() * this.emojiSymbols.length)];
  }

  private startAnimation() {
    // Change emojis every 10 seconds
    this.intervalId = setInterval(() => {
      this.emojis.forEach(emoji => {
        emoji.symbol = this.getRandomEmoji();
        // Occasionally change position slightly
        if (Math.random() > 0.7) {
          emoji.x = Math.max(0, Math.min(window.innerWidth - 100, emoji.x + (Math.random() - 0.5) * 200));
          emoji.y = Math.max(0, Math.min(window.innerHeight - 100, emoji.y + (Math.random() - 0.5) * 200));
        }
      });
    }, 10000);
  }

  trackByIndex(index: number): number {
    return index;
  }
}
