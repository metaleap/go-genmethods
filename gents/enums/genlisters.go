package gentenums

import (
	"strconv"

	. "github.com/go-leap/dev/go/gen"
	"github.com/go-leap/str"
	"github.com/metaleap/go-gent"
)

const (
	DefaultListBothFuncName     = "Wellknown{T}{s}"
	DefaultListBothDocComment   = "{N} returns the `names` and `values` of all {n} well-known `{T}` enumerants."
	DefaultListNamesFuncName    = "Wellknown{T}Names"
	DefaultListNamesDocComment  = "{N} returns the `names` of all {n} well-known `{T}` enumerants."
	DefaultListValuesFuncName   = "Wellknown{T}Values"
	DefaultListValuesDocComment = "{N} returns the `values` of all {n} well-known `{T}` enumerants."
	DefaultListMapFuncName      = "Wellknown{T}NamesAndValues"
	DefaultListMapDocComment    = "{N} returns the `namesToValues` of all {n} well-known `{T}` enumerants."
)

func init() {
	Gents.Listers.ListBoth.DocComment, Gents.Listers.ListBoth.Name = DefaultListBothDocComment, DefaultListBothFuncName
	Gents.Listers.ListNames.DocComment, Gents.Listers.ListNames.Name = DefaultListNamesDocComment, DefaultListNamesFuncName
	Gents.Listers.ListValues.DocComment, Gents.Listers.ListValues.Name = DefaultListValuesDocComment, DefaultListValuesFuncName
	Gents.Listers.ListMap.DocComment, Gents.Listers.ListMap.Name = DefaultListMapDocComment, DefaultListMapFuncName
}

// GentListEnumerantsFuncs generates a
// `func WellknownFoos() ([]string, []Foo)`
// for each enum type-def `Foo`.
//
// An instance with illustrative defaults is in `Gents.List`.
type GentListEnumerantsFuncs struct {
	gent.Opts

	ListNames  gent.Variation
	ListValues gent.Variation
	ListMap    gent.Variation
	ListBoth   gent.Variation

	AlwaysSkipFirst bool
	Rename          gent.Rename
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (me *GentListEnumerantsFuncs) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns) {
	if t.IsEnumish() && !(me.ListBoth.Disabled && me.ListMap.Disabled && me.ListNames.Disabled && me.ListValues.Disabled) {
		names, values := make(Syns, 0, len(t.Enumish.ConstNames)), make(Syns, 0, len(t.Enumish.ConstNames))
		for i, enumerant := range t.Enumish.ConstNames {
			if renamed := enumerant; enumerant != "_" && (i > 0 || !me.AlwaysSkipFirst) {
				if me.Rename != nil {
					renamed = me.Rename(ctx, t, enumerant)
				}
				if renamed != "" {
					names.Add(L(renamed))
					values.Add(t.G.T.N(enumerant))
				}
			}
		}

		var fnamevals, fnamenames string
		if !(me.ListBoth.Disabled && me.ListNames.Disabled && me.ListValues.Disabled) {
			fnamevals, fnamenames = me.ListValues.NameWith("T", t.Name), me.ListNames.NameWith("T", t.Name)
		}

		if !me.ListMap.Disabled {
			fnamemap := me.ListMap.NameWith("T", t.Name)
			yield.Add(me.genListMapFunc(t, fnamemap, names, values))
		}
		if !me.ListBoth.Disabled {
			fnameboth := me.ListBoth.NameWith("T", t.Name, "s", ustr.If(ustr.Suff(t.Name, "s"), "es", "s"))
			if ustr.Suff(fnameboth, "ys") {
				fnameboth = fnameboth[:len(fnameboth)-2] + "ies"
			}
			yield.Add(me.genListBothFunc(t, fnameboth, fnamenames, fnamevals, len(names)))
		}
		if !(me.ListNames.Disabled && me.ListBoth.Disabled) {
			yield.Add(me.genListNamesFunc(t, fnamenames, names))
		}
		if !(me.ListValues.Disabled && me.ListBoth.Disabled) {
			yield.Add(me.genListValuesFunc(t, fnamevals, values))
		}
	}
	return
}

func (me *GentListEnumerantsFuncs) genListNamesFunc(t *gent.Type, funcName string, enumerantNames Syns) *SynFunc {
	return Func(funcName).Ret("names", T.SliceOf.Strings).
		Doc(
			me.ListNames.DocComment.With("N", funcName, "T", t.Name, "n", strconv.Itoa(len(enumerantNames))),
		).
		Code(
			N("names").Set(L(enumerantNames)),
		)
}

func (me *GentListEnumerantsFuncs) genListValuesFunc(t *gent.Type, funcName string, enumerantValues Syns) *SynFunc {
	return Func(funcName).Ret("values", TSlice(t.G.T)).
		Doc(
			me.ListValues.DocComment.With("N", funcName, "T", t.Name, "n", strconv.Itoa(len(enumerantValues))),
		).
		Code(
			N("values").Set(L(enumerantValues)),
		)
}

func (me *GentListEnumerantsFuncs) genListBothFunc(t *gent.Type, funcName string, funcNameNames string, funcNameValues string, numEnumerants int) *SynFunc {
	return Func(funcName).Ret("names", T.SliceOf.Strings).Ret("values", TSlice(t.G.T)).
		Doc(
			me.ListBoth.DocComment.With("N", funcName, "T", t.Name, "n", strconv.Itoa(numEnumerants)),
		).
		Code(
			Names("names", "values").Set(Tup(Call(N(funcNameNames)), Call(N(funcNameValues)))),
		)
}

func (me *GentListEnumerantsFuncs) genListMapFunc(t *gent.Type, funcName string, enumerantNames Syns, enumerantValues Syns) *SynFunc {
	maptype := TMap(T.String, t.G.T)
	return Func(funcName).Ret("namesToValues", maptype).
		Doc(
			me.ListBoth.DocComment.With("N", funcName, "T", t.Name, "n", strconv.Itoa(len(enumerantNames))),
		).
		Code(
			N("namesToValues").Set(B.Make.Of(TMap(T.String, t.G.T), len(enumerantNames))),
			GEN_FOR(0, len(enumerantNames), 1, func(i int) ISyn {
				return N("namesToValues").At(enumerantNames[i]).Set(enumerantValues[i])
			}),
		)
}

// EnableOrDisableAllVariantsAndOptionals implements `github.com/metaleap/go-gent.IGent`.
func (me *GentListEnumerantsFuncs) EnableOrDisableAllVariantsAndOptionals(enabled bool) {
	disabled := !enabled
	me.ListBoth.Disabled, me.ListMap.Disabled, me.ListNames.Disabled, me.ListValues.Disabled = disabled, disabled, disabled, disabled
}
