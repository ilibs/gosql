## v2.0.0 (2019-09-12)
### Changed
- Remove the second argument to the Model() and Table() functions and replace it with WithTx(tx)
- Remove Model interface DbName() function,use the Use() function 
- Uniform API design specification, see [APIDESIGN](APIDESIGN.md)
- Relation add `connection:"db2"` struct tag, Solve the cross-library connection problem caused by deleting DbName()

## [v1.1.1 (2018-12-07)](https://github.com/ilibs/gosql/compare/v1.1.0...v1.1.1)

### Added
- Added Relation where

## [v1.1.0 (2018-12-06)](https://github.com/ilibs/gosql/compare/v1.0.10...v1.1.0)

### Added
- Added Relation

## [v1.0.10 (2018-12-03)](https://github.com/ilibs/gosql/compare/v1.0.9...v1.0.10)

### Added
- Added `gosql.Expr` Reference GORM Expr
- Added `In` Queries support