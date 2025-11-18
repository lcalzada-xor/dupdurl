package locale

import (
	"strings"
)

// TranslationGroup represents a group of translations for the same concept
type TranslationGroup struct {
	Canonical string   // Canonical form (usually English)
	Variants  []string // All known translations including canonical
}

// Common translations for typical web paths
var commonTranslations = []TranslationGroup{
	// About/Company pages
	{
		Canonical: "about",
		Variants: []string{
			"about", "about-us", "aboutus",
			"sobre-nosotros", "sobre", "acerca-de", "acerca", "quienes-somos", // Spanish
			"chi-siamo", "su-di-noi", "chi-sono", "riguardo", // Italian
			"a-propos", "qui-sommes-nous", // French
			"uber-uns", "ueber-uns", "wir", // German
			"sobre-nos", "quem-somos", // Portuguese
			"o-nas", "o-firme", // Polish/Czech
			"hakkimizda", "hakkinda", // Turkish
			"tentang-kami", "tentang", // Indonesian
		},
	},
	// Products/Services
	{
		Canonical: "products",
		Variants: []string{
			"products", "product",
			"productos", "producto", // Spanish
			"prodotti", "prodotto", // Italian
			"produits", "produit", // French
			"produkte", "produkt", // German
			"produtos", "produto", // Portuguese
			"produkty", "produkt", // Polish
			"urunler", "urun", // Turkish
		},
	},
	// Services
	{
		Canonical: "services",
		Variants: []string{
			"services", "service",
			"servicios", "servicio", // Spanish
			"servizi", "servizio", // Italian
			"services", "service", // French (same)
			"dienstleistungen", "dienste", // German
			"servicos", "servico", // Portuguese
			"uslugi", "usluga", // Polish/Russian
			"hizmetler", "hizmet", // Turkish
		},
	},
	// Contact
	{
		Canonical: "contact",
		Variants: []string{
			"contact", "contact-us", "contactus",
			"contacto", "contactanos", "contactenos", // Spanish
			"contatti", "contattaci", // Italian
			"contact", "contactez-nous", // French
			"kontakt", "kontaktieren", // German
			"contato", "fale-conosco", // Portuguese
			"kontakt", "kontaktuj", // Polish
			"iletisim", // Turkish
		},
	},
	// News/Blog
	{
		Canonical: "news",
		Variants: []string{
			"news", "blog", "articles",
			"noticias", "novedades", "articulos", // Spanish
			"notizie", "novita", "articoli", // Italian
			"nouvelles", "actualites", "blog", // French
			"nachrichten", "neuigkeiten", "blog", // German
			"noticias", "novidades", "artigos", // Portuguese
			"wiadomosci", "aktualnosci", // Polish
			"haberler", "blog", // Turkish
		},
	},
	// Help/Support
	{
		Canonical: "help",
		Variants: []string{
			"help", "support", "faq",
			"ayuda", "soporte", "preguntas-frecuentes", // Spanish
			"aiuto", "supporto", "domande-frequenti", // Italian
			"aide", "support", "faq", // French
			"hilfe", "support", "faq", // German
			"ajuda", "suporte", "perguntas-frequentes", // Portuguese
			"pomoc", "wsparcie", // Polish
			"yardim", "destek", // Turkish
		},
	},
	// Privacy/Legal
	{
		Canonical: "privacy",
		Variants: []string{
			"privacy", "privacy-policy",
			"privacidad", "politica-de-privacidad", // Spanish
			"privacy", "politica-sulla-privacy", // Italian
			"confidentialite", "politique-de-confidentialite", // French
			"datenschutz", "datenschutzrichtlinie", // German
			"privacidade", "politica-de-privacidade", // Portuguese
			"prywatnosc", "polityka-prywatnosci", // Polish
			"gizlilik", "gizlilik-politikasi", // Turkish
		},
	},
	{
		Canonical: "terms",
		Variants: []string{
			"terms", "terms-of-service", "terms-and-conditions",
			"terminos", "terminos-de-servicio", "condiciones", // Spanish
			"termini", "termini-di-servizio", "condizioni", // Italian
			"conditions", "conditions-utilisation", // French
			"bedingungen", "nutzungsbedingungen", "agb", // German
			"termos", "termos-de-servico", "condicoes", // Portuguese
			"warunki", "regulamin", // Polish
			"sartlar", "kullanim-kosullari", // Turkish
		},
	},
	// Account/User
	{
		Canonical: "account",
		Variants: []string{
			"account", "profile", "user",
			"cuenta", "perfil", "usuario", // Spanish
			"account", "profilo", "utente", // Italian
			"compte", "profil", "utilisateur", // French
			"konto", "profil", "benutzer", // German
			"conta", "perfil", "usuario", // Portuguese
			"konto", "profil", "uzytkownik", // Polish
			"hesap", "profil", "kullanici", // Turkish
		},
	},
	// Login/Signup
	{
		Canonical: "login",
		Variants: []string{
			"login", "signin", "sign-in",
			"iniciar-sesion", "ingresar", "entrar", // Spanish
			"accedi", "accesso", "login", // Italian
			"connexion", "se-connecter", // French
			"anmelden", "einloggen", "login", // German
			"entrar", "login", "iniciar-sessao", // Portuguese
			"zaloguj", "logowanie", // Polish
			"giris", "giris-yap", // Turkish
		},
	},
	{
		Canonical: "signup",
		Variants: []string{
			"signup", "register", "sign-up",
			"registrarse", "registro", "crear-cuenta", // Spanish
			"registrati", "registrazione", "iscriviti", // Italian
			"inscription", "sinscrire", "creer-compte", // French
			"registrieren", "anmelden", "konto-erstellen", // German
			"cadastro", "registrar", "criar-conta", // Portuguese
			"rejestracja", "zarejestruj", // Polish
			"kayit", "kayit-ol", "uye-ol", // Turkish
		},
	},
	// Home
	{
		Canonical: "home",
		Variants: []string{
			"home", "index", "main",
			"inicio", "principal", "casa", // Spanish
			"home", "inizio", "principale", // Italian
			"accueil", "index", "principale", // French
			"startseite", "home", "hauptseite", // German
			"inicio", "pagina-inicial", "principal", // Portuguese
			"strona-glowna", "start", // Polish
			"ana-sayfa", "anasayfa", "ev", // Turkish
		},
	},
	// Search
	{
		Canonical: "search",
		Variants: []string{
			"search", "find",
			"buscar", "busqueda", "encontrar", // Spanish
			"cerca", "ricerca", "trova", // Italian
			"recherche", "rechercher", "trouver", // French
			"suche", "suchen", "finden", // German
			"busca", "buscar", "procurar", // Portuguese
			"szukaj", "wyszukiwanie", // Polish
			"ara", "arama", "bul", // Turkish
		},
	},
	// Cart/Checkout
	{
		Canonical: "cart",
		Variants: []string{
			"cart", "basket", "shopping-cart",
			"carrito", "cesta", "canasta", // Spanish
			"carrello", "cestino", // Italian
			"panier", "chariot", // French
			"warenkorb", "einkaufswagen", // German
			"carrinho", "cesta", // Portuguese
			"koszyk", // Polish
			"sepet", "alisveris-sepeti", // Turkish
		},
	},
	{
		Canonical: "checkout",
		Variants: []string{
			"checkout", "payment", "pay",
			"pagar", "pago", "finalizar-compra", // Spanish
			"checkout", "pagamento", "paga", // Italian
			"paiement", "payer", "commander", // French
			"kasse", "bezahlen", "zahlung", // German
			"pagamento", "pagar", "finalizar", // Portuguese
			"kasa", "platnosc", // Polish
			"odeme", "odemeyap", // Turkish
		},
	},
}

// TranslationMatcher handles translation matching
type TranslationMatcher struct {
	normalizedIndex map[string]string // normalized variant -> canonical
	groupIndex      map[string]*TranslationGroup
}

// NewTranslationMatcher creates a new translation matcher
func NewTranslationMatcher() *TranslationMatcher {
	tm := &TranslationMatcher{
		normalizedIndex: make(map[string]string),
		groupIndex:      make(map[string]*TranslationGroup),
	}

	// Build indexes
	for i := range commonTranslations {
		group := &commonTranslations[i]
		canonical := normalizeForMatching(group.Canonical)

		tm.groupIndex[canonical] = group

		for _, variant := range group.Variants {
			normalized := normalizeForMatching(variant)
			tm.normalizedIndex[normalized] = canonical
		}
	}

	return tm
}

// AreTranslations checks if two path segments are translations of each other
func (tm *TranslationMatcher) AreTranslations(seg1, seg2 string) bool {
	norm1 := normalizeForMatching(seg1)
	norm2 := normalizeForMatching(seg2)

	// Same segment
	if norm1 == norm2 {
		return true
	}

	// Check if both belong to same translation group
	canonical1, ok1 := tm.normalizedIndex[norm1]
	canonical2, ok2 := tm.normalizedIndex[norm2]

	if ok1 && ok2 && canonical1 == canonical2 {
		return true
	}

	return false
}

// GetCanonical returns the canonical form of a segment if it's a known translation
func (tm *TranslationMatcher) GetCanonical(segment string) string {
	normalized := normalizeForMatching(segment)
	if canonical, ok := tm.normalizedIndex[normalized]; ok {
		return canonical
	}
	return segment
}

// normalizeForMatching normalizes a string for translation matching
func normalizeForMatching(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Remove common separators
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, "_", "")

	// Remove trailing 's' for simple pluralization
	if len(s) > 3 && strings.HasSuffix(s, "s") {
		return s[:len(s)-1]
	}

	return s
}
