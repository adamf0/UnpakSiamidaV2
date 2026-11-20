package event

import (
	"fmt"
	"html"
	"strings"
)

// RenderUserAccountEmailPlainText menghasilkan format teks sesuai spesifikasi
func RenderUserAccountEmailPlainText(r UserAccountEmailRecipient) string {
	roleDisplay := strings.TrimSpace(r.AbacRoles)
	if roleDisplay == "" {
		roleDisplay = r.Level
	}
	if roleDisplay == "" {
		roleDisplay = "user"
	}

	return fmt.Sprintf(`kami sampaikan akun SIAMIDA V3 sebagai berikut :

username : %s
password : %s
level : %s

buka siamida.unpak.ac.id dan pilih SIAMIDA V3.

Terimakasih`,
		r.Username,
		r.Password,
		roleDisplay,
	)
}

// RenderUserAccountEmailHTML merender template email profesional dengan layout modern responsif
func RenderUserAccountEmailHTML(r UserAccountEmailRecipient, loginURL string) string {
	if loginURL == "" {
		loginURL = "https://siamida.unpak.ac.id"
	}

	roleDisplay := strings.TrimSpace(r.AbacRoles)
	if roleDisplay == "" {
		roleDisplay = r.Level
	}
	if roleDisplay == "" {
		roleDisplay = "user"
	}

	nameDisplay := r.Name
	if strings.TrimSpace(nameDisplay) == "" {
		nameDisplay = r.Username
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Akun SIAMIDA V3 - Universitas Pakuan</title>
    <style>
        body {
            margin: 0;
            padding: 0;
            background-color: #f1f5f9;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
            color: #1e293b;
            line-height: 1.6;
        }
        .container {
            max-width: 600px;
            margin: 32px auto;
            background: #ffffff;
            border-radius: 16px;
            overflow: hidden;
            box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.05), 0 8px 10px -6px rgba(0, 0, 0, 0.05);
            border: 1px solid #e2e8f0;
        }
        .header {
            background: linear-gradient(135deg, #4f46e5 0%%, #7c3aed 100%%);
            padding: 36px 32px 32px;
            text-align: center;
            color: #ffffff;
        }
        .header-badge {
            display: inline-block;
            background: rgba(255, 255, 255, 0.2);
            backdrop-filter: blur(8px);
            padding: 4px 14px;
            border-radius: 9999px;
            font-size: 11px;
            font-weight: 700;
            letter-spacing: 1px;
            text-transform: uppercase;
            margin-bottom: 12px;
            border: 1px solid rgba(255, 255, 255, 0.25);
        }
        .header h1 {
            margin: 0 0 6px;
            font-size: 24px;
            font-weight: 800;
            letter-spacing: -0.5px;
        }
        .header p {
            margin: 0;
            font-size: 13px;
            color: #e0e7ff;
            opacity: 0.9;
        }
        .content {
            padding: 32px;
        }
        .greeting {
            font-size: 15px;
            margin-bottom: 16px;
            color: #334155;
        }
        .intro {
            font-size: 14px;
            color: #475569;
            margin-bottom: 24px;
        }
        .card-credentials {
            background: #f8fafc;
            border: 1px solid #e2e8f0;
            border-radius: 12px;
            padding: 20px 24px;
            margin-bottom: 28px;
        }
        .credential-row {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 10px 0;
            border-bottom: 1px solid #edf2f7;
        }
        .credential-row:last-child {
            border-bottom: none;
            padding-bottom: 0;
        }
        .credential-row:first-child {
            padding-top: 0;
        }
        .credential-label {
            font-size: 13px;
            font-weight: 600;
            color: #64748b;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }
        .credential-value {
            font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
            font-size: 14px;
            font-weight: 700;
            color: #0f172a;
            background: #ffffff;
            padding: 4px 10px;
            border-radius: 6px;
            border: 1px solid #cbd5e1;
        }
        .role-badge {
            background: #ecfdf5;
            color: #047857;
            border: 1px solid #a7f3d0;
            padding: 4px 12px;
            border-radius: 9999px;
            font-size: 12px;
            font-weight: 700;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }
        .instruction-box {
            background: #eff6ff;
            border-left: 4px solid #3b82f6;
            padding: 14px 18px;
            border-radius: 0 8px 8px 0;
            margin-bottom: 28px;
            font-size: 13.5px;
            color: #1e40af;
        }
        .action-container {
            text-align: center;
            margin: 32px 0;
        }
        .btn-portal {
            display: inline-block;
            background: linear-gradient(135deg, #4f46e5 0%%, #4338ca 100%%);
            color: #ffffff !important;
            text-decoration: none;
            font-weight: 700;
            font-size: 14px;
            padding: 14px 32px;
            border-radius: 10px;
            box-shadow: 0 4px 14px rgba(79, 70, 229, 0.35);
            letter-spacing: 0.3px;
        }
        .security-note {
            font-size: 12px;
            color: #64748b;
            background: #fffbeb;
            border: 1px solid #fef3c7;
            padding: 12px 16px;
            border-radius: 8px;
            margin-top: 24px;
            line-height: 1.5;
        }
        .footer {
            padding: 24px 32px;
            background: #f8fafc;
            border-top: 1px solid #e2e8f0;
            font-size: 13px;
            color: #475569;
        }
        .footer strong {
            color: #1e293b;
        }
        .footer-sub {
            margin-top: 14px;
            font-size: 11px;
            color: #94a3b8;
            border-top: 1px solid #e2e8f0;
            padding-top: 12px;
            text-align: center;
        }
    </style>
</head>
<body>
    <div class="container">
        <!-- HEADER -->
        <div class="header">
            <div class="header-badge">Universitas Pakuan</div>
            <h1>SIAMIDA V3</h1>
            <p>Sistem Penjaminan Mutu &amp; Akuntabilitas Kinerja</p>
        </div>

        <!-- CONTENT -->
        <div class="content">
            <div class="greeting">
                Yth. Bapak/Ibu <strong>%s</strong>,
            </div>
            <div class="intro">
                Bersama ini kami sampaikan akun sistem informasi <strong>SIAMIDA V3</strong> sebagai berikut:
            </div>

            <!-- CREDENTIALS CARD -->
            <table class="card-credentials" width="100%%" cellpadding="0" cellspacing="0" style="background:#f8fafc; border:1px solid #e2e8f0; border-radius:12px; margin-bottom:24px;">
                <tr>
                    <td style="padding:14px 20px; border-bottom:1px solid #e2e8f0; font-size:13px; font-weight:600; color:#64748b; width:35%%;">
                        USERNAME
                    </td>
                    <td style="padding:14px 20px; border-bottom:1px solid #e2e8f0; text-align:right;">
                        <span style="font-family:monospace; font-size:14px; font-weight:700; color:#0f172a; background:#ffffff; padding:4px 10px; border-radius:6px; border:1px solid #cbd5e1;">
                            %s
                        </span>
                    </td>
                </tr>
                <tr>
                    <td style="padding:14px 20px; border-bottom:1px solid #e2e8f0; font-size:13px; font-weight:600; color:#64748b;">
                        PASSWORD
                    </td>
                    <td style="padding:14px 20px; border-bottom:1px solid #e2e8f0; text-align:right;">
                        <span style="font-family:monospace; font-size:14px; font-weight:700; color:#4f46e5; background:#ffffff; padding:4px 10px; border-radius:6px; border:1px solid #cbd5e1;">
                            %s
                        </span>
                    </td>
                </tr>
                <tr>
                    <td style="padding:14px 20px; font-size:13px; font-weight:600; color:#64748b;">
                        LEVEL / ROLE
                    </td>
                    <td style="padding:14px 20px; text-align:right;">
                        <span style="background:#ecfdf5; color:#047857; border:1px solid #a7f3d0; padding:4px 12px; border-radius:9999px; font-size:12px; font-weight:700; text-transform:uppercase;">
                            %s
                        </span>
                    </td>
                </tr>
            </table>

            <!-- INSTRUCTION CALLOUT -->
            <div style="background:#eff6ff; border-left:4px solid #3b82f6; padding:14px 18px; border-radius:0 8px 8px 0; margin-bottom:24px; font-size:13.5px; color:#1e40af;">
                Silakan buka <strong>siamida.unpak.ac.id</strong> dan pilih <strong>SIAMIDA V3</strong> untuk masuk ke dalam sistem.
            </div>

            <!-- BUTTON CTA -->
            <div style="text-align:center; margin:28px 0;">
                <a href="%s" target="_blank" style="display:inline-block; background:linear-gradient(135deg, #4f46e5 0%%, #4338ca 100%%); color:#ffffff; text-decoration:none; font-weight:700; font-size:14px; padding:14px 32px; border-radius:10px; box-shadow:0 4px 14px rgba(79, 70, 229, 0.35);">
                    Buka SIAMIDA V3 &rarr;
                </a>
            </div>

            <!-- SECURITY NOTICE -->
            <div style="font-size:12px; color:#92400e; background:#fef3c7; border:1px solid #fde68a; padding:12px 16px; border-radius:8px; line-height:1.5;">
                <strong>Perhatian Keamanan:</strong> Demi keamanan akun Anda, harap segera login dan perbarui kata sandi Anda. Jangan berikan kredensial ini kepada pihak yang tidak berkepentingan.
            </div>
        </div>

        <!-- FOOTER -->
        <div class="footer">
            <p style="margin:0 0 4px;">Terimakasih,</p>
            <p style="margin:0; font-weight:700; color:#1e293b;">
                Lembaga Penjaminan Mutu (LPM)<br>
                Universitas Pakuan
            </p>
        </div>
    </div>
</body>
</html>`,
		html.EscapeString(nameDisplay),
		html.EscapeString(r.Username),
		html.EscapeString(r.Password),
		html.EscapeString(roleDisplay),
		html.EscapeString(loginURL),
	)
}

// RenderUserAccountEmailTelegram menghasilkan notifikasi telegram untuk audit
func RenderUserAccountEmailTelegram(event UserSendAccountEmailEvent) string {
	var recLines []string
	for _, u := range event.Users {
		roleDisplay := u.AbacRoles
		if roleDisplay == "" {
			roleDisplay = u.Level
		}
		recLines = append(recLines, fmt.Sprintf("• <b>%s</b> (%s) - Role: <code>%s</code>",
			html.EscapeString(u.Username),
			html.EscapeString(u.Email),
			html.EscapeString(roleDisplay),
		))
	}

	return fmt.Sprintf(
		"<b>[OUTBOX EMAIL AKUN SIAMIDA V3]</b>\nTotal Terkirim: <b>%d akun</b>\nDaftar Penerima:\n%s",
		len(event.Users),
		strings.Join(recLines, "\n"),
	)
}
